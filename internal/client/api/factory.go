package api

import (
	"github.com/ncyellow/GophKeeper/internal/client/config"
)

// CreateSender function creates either an https or grpc client based on the settings
func CreateSender(conf *config.Config) (Sender, error) {
	// By default, we use https, only if the GRPCAddress entrypoint is specified, we switch to grpc
	if conf.GRPCAddress != "" {
		// establish a connection to the server
		return NewGRPCSender(conf)
	}
	return NewHTTPSender(conf)
}
