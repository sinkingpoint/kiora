package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/sinkingpoint/kiora/lib/kiora/model"
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

	// The port that the cluster communication end of this instance is attached to.
	clusterPort string

	shutdownOnce sync.Once
	shutdown     bool
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

func (k *KioraInstance) WithName(name string) *KioraInstance {
	k.name = name
	return k
}

func (k *KioraInstance) WithConfigFile(configFile string) *KioraInstance {
	k.configFile = configFile
	return k
}

// Start actually executes the Kiora command, running it in a background go routine.
func (k *KioraInstance) Start(t *testing.T) error {
	t.Helper()
	name := kioraInstanceName()
	if k.name == "" {
		k.name = name
	}

	httpPort, err := getRandomPort()
	require.NoError(t, err)

	clusterPort, err := getRandomPort()
	require.NoError(t, err)

	args := append([]string{"run", "../cmd/kiora", "-c", k.configFile, "--web.listen-url", "localhost:" + httpPort, "--cluster.listen-url", "localhost:" + clusterPort}, k.args...)

	k.httpPort = httpPort
	k.clusterPort = clusterPort
	k.cmd = exec.Command("go", args...)
	k.cmd.Stdout = k.stdout
	k.cmd.Stderr = k.stderr

	// Set up a dedicated process group, so we can kill every child process.
	k.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Setup a cleanup job that stops the instance, and removes the data directory.
	t.Cleanup(func() {
		t.Logf("Name: %q Stderr: \n%s", k.name, k.Stderr())
		t.Logf("Name: %q Stdout: \n%s", k.name, k.Stdout())
		if !k.shutdown {
			require.NoError(t, k.Stop())
		}
		require.NoError(t, os.RemoveAll("../artifacts/test/"+name))
	})

	go func() {
		k.exitChannel <- k.cmd.Run()
	}()

	return k.WaitTillUp(context.TODO(), t)
}

func (k *KioraInstance) IsUp(ctx context.Context, t *testing.T) bool {
	url := k.GetHTTPURL("/")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err)

	_, err = http.DefaultClient.Do(req)
	return err == nil
}

func (k *KioraInstance) WaitTillUp(ctx context.Context, t *testing.T) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			if k.IsUp(ctx, t) {
				return nil
			}
		case <-ctx.Done():
			return errors.New("didn't come up within context")
		}
	}
}

// GetURL returns a call to this instance, on the given path. This ellides the need to interact with the ports on this instance directly.
func (k *KioraInstance) GetHTTPURL(path string) string {
	return "http://localhost:" + k.httpPort + path
}

func (k *KioraInstance) GetClusterHost() string {
	return "localhost:" + k.clusterPort
}

// Stop sends a sigkill to the process group that backs this instance.
func (k *KioraInstance) Stop() error {
	var err error
	k.shutdownOnce.Do(func() {
		err = syscall.Kill(-k.cmd.Process.Pid, syscall.SIGKILL)
		k.shutdown = true
	})

	return err
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
	requestURL := k.GetHTTPURL("/api/v1/alerts")

	alertBytes, err := json.Marshal([]model.Alert{alert})
	require.NoError(t, err)

	resp, err := http.Post(requestURL, "application/json", bytes.NewReader(alertBytes))
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, http.StatusAccepted, resp.StatusCode, "body: %s", string(body))
}

// kioraInstanceName returns a 16 character long random string that will be used as the name of a KioraInstance.
func kioraInstanceName() string {
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

// StartKioraCluster starts n KioraInstance's, binding them into a serf cluster.
func StartKioraCluster(t *testing.T, numNodes int) []*KioraInstance {
	t.Helper()

	// Start a leader node, telling it to bootstrap the cluster.
	leader := NewKioraInstance().WithName("node-0")
	require.NoError(t, leader.Start(t))

	// Start n-1 instances.
	nodes := []*KioraInstance{}
	for i := 1; i < numNodes; i++ {
		node := NewKioraInstance("--cluster.bootstrap-peers", leader.GetClusterHost()).WithName(fmt.Sprintf("node-%d", i))
		require.NoError(t, node.Start(t))
		nodes = append(nodes, node)
	}

	nodes = append(nodes, leader)

	// Wait for a bit, to let the gossip to settle.
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
