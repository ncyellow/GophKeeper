// Package gprcserver содержит имплементацию сервера через grpc
package gprcserver

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/gprcserver/api"
	"github.com/ncyellow/GophKeeper/internal/server/gprcserver/proto"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// GRPCServer структура сервера. Реализует интерфейс Server
type GRPCServer struct {
	Conf *config.Config
}

// Run блокирующая функция запуска сервера.
// После запуска встает в ожидание os.Interrupt, syscall.SIGINT, syscall.SIGTERM
// Функция очень похожа на RunServer из http реализации, но тут другой вариант graceful shutdown.
func (s *GRPCServer) Run() {
	store, err := storage.NewPgStorage(s.Conf)
	if err != nil {
		log.Fatal().Err(err)
	}

	listen, err := net.Listen("tcp", s.Conf.GRPCAddress)
	if err != nil {
		log.Fatal().Err(err)
	}

	grpcServer := grpc.NewServer()
	// регистрируем сервис
	proto.RegisterGophKeeperServerServer(grpcServer, api.NewServer(store, s.Conf))

	defer func() {
		// гасим сервер через GracefulStop
		grpcServer.GracefulStop()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Error().Err(err)
		}
	}()

	<-done
	log.Info().Msg("Server Shutdown gracefully")
}
