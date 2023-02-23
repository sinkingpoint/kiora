package clustering

import (
	"context"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

const DEFAULT_REFRESH_TIME = 2 * time.Second

// ObserverID is an arbitrary number used for removing observers.
type ObserverID uint32

// Observer receives notifications when the servers in a Clusterer change.
type Observer interface {
	AddServer(server Server)
	RemoveServer(server Server)
}

// StateObserver periodically polls a Clusterer, comparing the servers
// to the previous and updating Observers when the state changes.
type StateObserver struct {
	refreshTime     time.Duration
	backend         Clusterer
	previousServers map[string]Server
	observers       map[ObserverID]Observer
	killChan        chan struct{}
}

func NewStateObserver(backend Clusterer) *StateObserver {
	return &StateObserver{
		refreshTime:     DEFAULT_REFRESH_TIME,
		backend:         backend,
		killChan:        make(chan struct{}),
		previousServers: make(map[string]Server),
		observers:       make(map[ObserverID]Observer),
	}
}

func (s *StateObserver) WithRefreshInterval(interval time.Duration) *StateObserver {
	s.refreshTime = interval
	return s
}

// AddObserver adds the given observer so that it receives state updates.
func (s *StateObserver) AddObserver(o Observer) ObserverID {
	id := ObserverID(rand.Uint32())
	s.observers[id] = o

	return id
}

// RemoveObserver removes the observer with the given ID so that it no longer receives updates.
func (s *StateObserver) RemoveObserver(id ObserverID) {
	delete(s.observers, id)
}

// removeServer communicates the server to all the configured observers.
func (s *StateObserver) removeServer(server Server) {
	for _, o := range s.observers {
		o.RemoveServer(server)
	}

	delete(s.previousServers, server.Name())
}

// addServer communicates the server to all the configured observers.
func (s *StateObserver) addServer(server Server) {
	for _, o := range s.observers {
		o.AddServer(server)
	}

	s.previousServers[server.Name()] = server
}

// observe polls the clusterer and calls addServer and removeServer as necessary
// to reconcile the old state with the new one.
func (s *StateObserver) observe() {
	ctx := context.Background()
	servers, err := s.backend.GetMembers(ctx)
	if err != nil {
		log.Err(err).Msg("failed to retrieve cluster state")
		return
	}

	currentServers := map[string]Server{}
	for _, s := range servers {
		currentServers[s.Name()] = s
	}

	for name, server := range s.previousServers {
		if _, ok := currentServers[name]; !ok {
			s.removeServer(server)
		}
	}

	for name, server := range currentServers {
		if _, ok := s.previousServers[name]; !ok {
			s.addServer(server)
		}
	}

	// TODO(cdouch): Would it be useful to have an `updateServer` for when addresses change?
}

// Run starts the loop that runs this StateObserver, until Kill is called. Note: This blocks, so run it in a goroutine.
func (s *StateObserver) Run() {
	ticker := time.NewTicker(s.refreshTime)
	for {
		select {
		case <-ticker.C:
			s.observe()
		case <-s.killChan:
			return
		}
	}
}

// Kill stops this state observer after a call to `Run`.
func (s *StateObserver) Kill() {
	s.killChan <- struct{}{}
}
