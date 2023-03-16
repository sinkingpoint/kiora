package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Service is a service that runs in the background of Kiora, performing some task.
type Service interface {
	// Name returns the human readable name of the service.
	Name() string

	// Run runs the service until the given context is done.
	Run(ctx context.Context) error
}

// BackgroundServices wraps a number of services and provides a way to cancel all of them.
type BackgroundServices struct {
	services []Service

	ctx    context.Context
	cancel context.CancelFunc
}

func NewBackgroundServices() *BackgroundServices {
	return &BackgroundServices{}
}

func (b *BackgroundServices) RegisterService(s Service) {
	b.services = append(b.services, s)
}

func (b *BackgroundServices) Run(ctx context.Context) error {
	b.ctx, b.cancel = context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	var err error
	for _, s := range b.services {
		wg.Add(1)
		go func(s Service) {
			if err = s.Run(b.ctx); err != nil {
				err = errors.Wrapf(err, "service %q failed", s.Name())
			}

			// The underlying context is still open, but the service has exitted. Stop the world.
			if b.ctx.Err() == nil {
				b.Shutdown(b.ctx)
				err = fmt.Errorf("service %q failed without an error", s.Name())
			}

			wg.Done()
		}(s)
	}

	wg.Wait()

	log.Info().Msg("Background Services Shut Down")

	return err
}

func (b *BackgroundServices) Shutdown(ctx context.Context) {
	b.cancel()
}
