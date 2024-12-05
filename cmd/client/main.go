// Package main contains the client launch.
package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ncyellow/GophKeeper/internal/client/api"
	"github.com/ncyellow/GophKeeper/internal/client/config"
	"github.com/ncyellow/GophKeeper/internal/client/console"
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
