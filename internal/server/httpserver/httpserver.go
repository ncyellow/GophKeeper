// Package httpserver - реализует создание и запуск сервера
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

// HTTPServer структура нашего https сервера. Реализует интерфейс Server
type HTTPServer struct {
	Conf *config.Config
}

// Run блокирующая функция по запуску сервера
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
		if err := srv.ListenAndServeTLS(s.Conf.CryptoCrt, s.Conf.CryptoKey); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Msgf("listen: %s", err)
		}
	}()
	<-idleConnsClosed
	log.Info().Msg("Server Shutdown gracefully")
	return nil
}
