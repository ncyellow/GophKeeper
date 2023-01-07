package api

import (
	"github.com/ncyellow/GophKeeper/internal/client/config"
)

func CreateSender(conf *config.Config) Sender {
	// По дефолту у нас http, только если задан GRPCAddress entrypoint, мы переходим на grpc
	if conf.GRPCAddress != "" {
		// устанавливаем соединение с сервером
		return NewGRPCSender(conf)
	}
	return NewHTTPSender(conf)
}
