package api

import (
	"github.com/ncyellow/GophKeeper/internal/client/config"
)

// CreateSender функция создает либо клиент https либо grpc по настройкам
func CreateSender(conf *config.Config) (Sender, error) {
	// По дефолту у нас https, только если задан GRPCAddress entrypoint, мы переходим на grpc
	if conf.GRPCAddress != "" {
		// устанавливаем соединение с сервером
		return NewGRPCSender(conf)
	}
	return NewHTTPSender(conf)
}
