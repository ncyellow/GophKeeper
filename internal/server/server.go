// Package server содержит реализацию сервера нашего хранилища данных. Содержит дле реализации grpc и https
package server

import (
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/gprcserver"
	"github.com/ncyellow/GophKeeper/internal/server/httpserver"
)

// Server интерфейс который должен реализовать сервер
type Server interface {
	Run()
}

// CreateServer - factory function которая выбирает имплементацию сервера по параметрам
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
