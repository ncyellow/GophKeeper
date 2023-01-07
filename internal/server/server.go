package server

import (
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/gprcserver"
	"github.com/ncyellow/GophKeeper/internal/server/httpserver"
)

type Server interface {
	Run()
}

func CreateServer(conf *config.Config) Server {
	// По дефолту у нас http, только если задан GRPCAddress entrypoint, мы переходим на grpc
	if conf.GRPCAddress != "" {
		// устанавливаем соединение с сервером
		return &gprcserver.GRPCServer{
			Conf: config.ParseConfig(),
		}
	}
	return &httpserver.HTTPServer{
		Conf: config.ParseConfig(),
	}
}
