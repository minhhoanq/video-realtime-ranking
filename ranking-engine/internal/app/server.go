package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"video-realtime-ranking/ranking-engine/pkg/constants"

	"golang.org/x/sync/errgroup"
)

const (
	_defaultShutdownTimeout = constants.DefaultShutdownTimeout
	_defaultAddr            = constants.DefaultPort
	_defaultReadTimeout     = constants.DefaultReadTimeout
	_defaultWriteTimeout    = constants.DefaultWriteTimeout
)

type Server struct {
	server    *http.Server
	waitGroup *errgroup.Group
	notify    chan error
	ctx       context.Context
}

func NewServer(handler http.Handler, waitGroup *errgroup.Group, ctx context.Context, opts ...Option) {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &Server{
		server:    httpServer,
		waitGroup: waitGroup,
		notify:    make(chan error, 1),
		ctx:       ctx,
	}

	// custom options
	for _, opt := range opts {
		opt(s)
	}

	s.Start()
	s.Shutdown()
}

func (s *Server) Start() {
	s.waitGroup.Go(func() error {
		fmt.Println("start HTTP server at Address: ", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		}

		return nil
	})
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() {
	s.waitGroup.Go(func() error {
		<-s.ctx.Done()
		err := s.server.Shutdown(context.Background())
		if err != nil {
			fmt.Println("failed to shutdown HTTP server error: ", err.Error())
		}
		fmt.Println("graceful shutdown HTTP server")
		fmt.Println("HTTP server is stopped")
		return nil
	})
}
