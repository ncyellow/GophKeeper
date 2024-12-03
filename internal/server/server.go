// Package server contains the implementation of our data storage server. It includes implementations for both grpc and https
package server

import (
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/gprcserver"
	"github.com/ncyellow/GophKeeper/internal/server/httpserver"
)

// Server interface that the server must implement
type Server interface {
	Run() error
}

// CreateServer - factory function that selects the server implementation based on parameters
func CreateServer(conf *config.Config) Server {
	// By default, we use http, only if the GRPCAddress entrypoint is specified, we switch to grpc
	if conf.GRPCAddress != "" {
		// establish a connection to the server
		return &gprcserver.GRPCServer{
			Conf: conf,
		}
	}
	return &httpserver.HTTPServer{
		Conf: conf,
	}
}
