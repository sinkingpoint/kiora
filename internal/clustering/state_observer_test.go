package clustering_test

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/mocks/mock_clustering"
)

type testServer struct {
	name    string
	address string
}

func NewTestServer(name string, address string) clustering.Server {
	return testServer{
		name:    name,
		address: address,
	}
}

func (t testServer) Name() string {
	return t.name
}

func (t testServer) Address() string {
	return t.address
}

func TestStateObserver(t *testing.T) {
	ctrl := gomock.NewController(t)

	clusterer := mock_clustering.NewMockClusterer(ctrl)
	clusterer.EXPECT().GetMembers(gomock.Any()).Return([]clustering.Server{
		NewTestServer("foo", "foo"),
	}, nil)

	clusterer.EXPECT().GetMembers(gomock.Any()).Return([]clustering.Server{
		NewTestServer("bar", "bar"),
	}, nil)

	// The observer should see three requests - adding the first server, removing it, and adding the second.
	observer := mock_clustering.NewMockObserver(ctrl)
	observer.EXPECT().AddServer(NewTestServer("foo", "foo")).Times(1)
	observer.EXPECT().RemoveServer(NewTestServer("foo", "foo")).Times(1)
	observer.EXPECT().AddServer(NewTestServer("bar", "bar")).Times(1)

	stateObserver := clustering.NewStateObserver(clusterer).WithRefreshInterval(500 * time.Millisecond)
	stateObserver.AddObserver(observer)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		stateObserver.Run()
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	stateObserver.Kill()
	wg.Wait()
}
