// Package gprcserver contains the implementation of the server via gRPC
package gprcserver

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ncyellow/GophKeeper/internal/proto"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/gprcserver/api"
	"github.com/ncyellow/GophKeeper/internal/server/storage"
)

// GRPCServer server structure. Implements the Server interface
type GRPCServer struct {
	Conf *config.Config
}

// Run blocking function for server startup.
// After startup, it waits for os.Interrupt, syscall.SIGINT, syscall.SIGTERM
// The function is very similar to RunServer from the http implementation, but here is a different variant of graceful shutdown.
func (s *GRPCServer) Run() error {
	store, err := storage.NewPgStorage(s.Conf)
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", s.Conf.GRPCAddress)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	// register service
	proto.RegisterGophKeeperServerServer(grpcServer, api.NewServer(store, s.Conf))

	defer func() {
		// shutting down the server via GracefulStop
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

	return nil
}
