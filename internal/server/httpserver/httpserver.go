// Package httpserver implements the creation and launch of a server
package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/ncyellow/GophKeeper/internal/server/auth/jwt"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
)

// HTTPServer structure of our HTTPS server. Implements the Server interface
type HTTPServer struct {
	Conf *config.Config
}

// Run a blocking function for starting the server
func (s *HTTPServer) Run() error {
	store, err := storage.NewPgStorage(s.Conf)
	if err != nil {
		return err
	}

	router := NewRouter(s.Conf, store, &jwt.DefaultParser{})

	srv := http.Server{
		Addr:    s.Conf.Address,
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// waiting for an interrupt
		<-done
		// shutting down the server
		if err := srv.Shutdown(context.Background()); err != nil {
			// errors closing the Listener
			log.Info().Msgf("HTTP server Shutdown: %v", err)
		}
		// notifying the main thread
		// that all network connections have been processed and closed
		close(idleConnsClosed)
	}()

	go func() {
		if err := srv.ListenAndServeTLS(s.Conf.CryptoCrt, s.Conf.CryptoKey); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Msgf("listen: %s", err)
		}
	}()
	<-idleConnsClosed
	log.Info().Msg("Server Shutdown gracefully")
	return nil
}
