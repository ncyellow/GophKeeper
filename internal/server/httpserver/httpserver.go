// Package httpserver - реализует создание и запуск сервера
package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
	"github.com/rs/zerolog/log"
)

type HTTPServer struct {
	Conf *config.Config
}

func (s *HTTPServer) Run() {
	store := storage.NewPgStorage(s.Conf)
	router := NewRouter(s.Conf, store)

	srv := http.Server{
		Addr:    s.Conf.Address,
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// ждем прерывание
		<-done
		// гасим сервер
		if err := srv.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			log.Info().Msgf("HTTP server Shutdown: %v", err)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	go func() {
		if err := http.ListenAndServe(s.Conf.Address, router); err != nil && err != http.ErrServerClosed {
			log.Error().Msgf("listen: %s", err)
		}
	}()
	<-idleConnsClosed
	log.Info().Msg("Server Shutdown gracefully")
}
