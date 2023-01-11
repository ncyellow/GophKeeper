// Package main содержит запуск клиента.
package main

import (
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/client/api"
	"github.com/ncyellow/GophKeeper/internal/client/config"
	"github.com/ncyellow/GophKeeper/internal/client/console"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("Client run")

	conf, err := config.ParseConfig()
	if err != nil {
		log.Fatal().Err(err)
	}

	sender, err := api.CreateSender(conf)
	if err != nil {
		log.Fatal().Err(err)
	}

	terminal := console.Console{
		Conf:   conf,
		Client: sender,
	}
	terminal.Run()
}
