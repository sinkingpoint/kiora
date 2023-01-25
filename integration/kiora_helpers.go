package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// KioraInstance wraps an instance of Kiora started as a seperate process as a black box.
type KioraInstance struct {
	// The cluster name of this instance.
	name string

	// The path to the file that contains the config for this instance.
	configFile string

	// The CLI arguments to add to the command to start Kiora.
	args []string

	// The channel that contains the error from `cmd.Run`
	exitChannel chan error

	// The backing command of this instance.
	cmd *exec.Cmd

	// The stdout stream of this instance.
	stdout *bytes.Buffer

	// The stderr stream of this instance.
	stderr *bytes.Buffer

	// The port that the HTTP end of this instance is attached to.
	httpPort string

	// The port that the raft end of this instance is attached to.
	raftPort string
}

// NewKioraInstance constructs a new KioraInstance that will start a Kiora run with the given CLI args.
func NewKioraInstance(args ...string) *KioraInstance {
	return &KioraInstance{
		args:        args,
		exitChannel: make(chan error),
		configFile:  "../testdata/kiora.dot",
		stdout:      &bytes.Buffer{},
		stderr:      &bytes.Buffer{},
	}
}

func (k *KioraInstance) WithConfigFile(configFile string) *KioraInstance {
	k.configFile = configFile
	return k
}

// Start actually executes the Kiora command, running it in a background go routine.
func (k *KioraInstance) Start(t *testing.T) error {
	t.Helper()
	name := kioraInstanceName()
	httpPort, err := getRandomPort()
	require.NoError(t, err)

	raftPort, err := getRandomPort()
	require.NoError(t, err)

	args := append([]string{"run", "../cmd/kiora", "-c", k.configFile, "--raft.data-dir",
		"../artifacts/test/" + name, "--web.listen-url", "localhost:" + httpPort,
		"--raft.listen-url", "localhost:" + raftPort}, k.args...)

	k.name = name
	k.httpPort = httpPort
	k.raftPort = raftPort
	k.cmd = exec.Command("go", args...)
	k.cmd.Stdout = k.stdout
	k.cmd.Stderr = k.stderr

	// Set up a dedicated process group, so we can kill every child process.
	k.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Setup a cleanup job that stops the instance, and removes the data directory.
	t.Cleanup(func() {
		t.Logf("Name: %q Stderr: \n%s", k.args, k.Stderr())
		t.Logf("Name: %q Stdout: \n%s", k.args, k.Stdout())
		require.NoError(t, k.Stop())
		require.NoError(t, os.RemoveAll("../artifacts/test/"+name))
	})

	go func() {
		k.exitChannel <- k.cmd.Run()
	}()

	return nil
}

// clusterHasLeader checks the raft cluster status, and returns true if any node in the cluster is the leader.
func (k *KioraInstance) clusterHasLeader(t *testing.T) bool {
	reqURL := k.GetURL("/admin/raft/status")
	resp, err := http.Get(reqURL)
	if err != nil {
		return false
	}

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	resp.Body.Close()

	return strings.Contains(string(body), `"is_leader":true`)
}

// WaitUntilLeader polls the raft endpoint until the cluster has a leader, failing if it isn't up within 10 seconds.
func (k *KioraInstance) WaitUntilLeader(t *testing.T, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if k.clusterHasLeader(t) {
				return nil
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

// JoinWith adds in the given KioraInstance to the cluster that contains `k`
func (k *KioraInstance) JoinWith(k2 *KioraInstance) error {
	reqURL := k.GetURL("/admin/raft/add_member")
	resp, err := http.Post(reqURL, "application/json", strings.NewReader(fmt.Sprintf(`{"id":"%s","address":"%s"}`, k2.name, "localhost:"+k2.raftPort)))
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add to raft cluster")
	}

	return err
}

// GetURL returns a call to this instance, on the given path. This ellides the need to interact with the ports on this instance directly.
func (k *KioraInstance) GetURL(path string) string {
	return "http://localhost:" + k.httpPort + path
}

// Stop sends a sigkill to the process group that backs this instance.
func (k *KioraInstance) Stop() error {
	return syscall.Kill(-k.cmd.Process.Pid, syscall.SIGKILL)
}

// Stdout returns the contents of the stdout stream of this instance.
func (k *KioraInstance) Stdout() string {
	return k.stdout.String()
}

// Stderr returns the contents of the stderr stream of this instance.
func (k *KioraInstance) Stderr() string {
	return k.stderr.String()
}

// WaitForExit waits for the instance to finish, returning an error if the context expires.
func (k *KioraInstance) WaitForExit(ctx context.Context) error {
	select {
	case err := <-k.exitChannel:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (k *KioraInstance) SendAlert(t *testing.T, ctx context.Context, alert model.Alert) {
	requestURL := k.GetURL("/api/v1/alerts")

	alertBytes, err := json.Marshal([]model.Alert{alert})
	require.NoError(t, err)

	t.Log(string(alertBytes))

	resp, err := http.Post(requestURL, "application/json", bytes.NewReader(alertBytes))
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, http.StatusAccepted, resp.StatusCode)
}

// ClusterStatus is a convenience method that returns the output of the raft cluster status endpoint.
func (k *KioraInstance) ClusterStatus(t *testing.T) string {
	t.Helper()
	statusResp, err := http.Get(k.GetURL("/admin/raft/status"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusResp.StatusCode)

	s, err := io.ReadAll(statusResp.Body)
	assert.NoError(t, err)
	statusResp.Body.Close()

	return string(s)
}

// kioraInstanceName returns a 16 character long random string that will be used as the name of a KioraInstance.
func kioraInstanceName() string {
	rand.Seed(time.Now().UnixNano())
	n := 16
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// getRandomPort attempts to get a random, unused port.
func getRandomPort() (string, error) {
	// TL;DR, let the OS allocate us a random port, stop listening on it, and then return it, hoping the port isn't used when we attempt to rebind.
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	url, err := url.Parse("http://" + l.Addr().String())
	if err != nil {
		return "", err
	}

	if err := l.Close(); err != nil {
		return "", err
	}

	return url.Port(), nil
}

// StartKioraCluster starts n KioraInstance's, binding them into a raft cluster.
func StartKioraCluster(t *testing.T, numNodes int) []*KioraInstance {
	t.Helper()
	// Start n-1 instances.
	nodes := []*KioraInstance{}
	for i := 0; i < numNodes-1; i++ {
		node := NewKioraInstance("--raft.local-id", fmt.Sprintf("node-%d", i))
		require.NoError(t, node.Start(t))
		nodes = append(nodes, node)
	}

	// Start a leader node, telling it to bootstrap the cluster.
	leader := NewKioraInstance("--raft.bootstrap", "--raft.local-id", fmt.Sprintf("node-%d", numNodes-1))
	require.NoError(t, leader.Start(t))

	// Wait until the cluster is up, and then add every node to the leader.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, leader.WaitUntilLeader(t, ctx))
	for _, node := range nodes {
		require.NoError(t, leader.JoinWith(node))
	}

	nodes = append(nodes, leader)

	// Wait for a bit, to let the raft cluster settle.
	time.Sleep(2 * time.Second)

	return nodes
}

// WriteConfigFile writes out the given config to a file, returning
// the path that can be added to a Kiora instance.
func WriteConfigFile(t *testing.T, config string) string {
	t.Helper()
	file, err := os.CreateTemp("", "")
	require.NoError(t, err)

	n, err := file.WriteString(config)
	require.NoError(t, err)
	require.Equal(t, len(config), n)

	t.Cleanup(func() {
		os.Remove(file.Name())
	})

	return file.Name()
}
