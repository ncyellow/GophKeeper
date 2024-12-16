// Package main contains the server launch.
package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ncyellow/GophKeeper/internal/server"
	"github.com/ncyellow/GophKeeper/internal/server/config"
)

func main() {
	fmt.Println("Server start")

	conf, err := config.ParseConfig()
	if err != nil {
		log.Fatal().Err(err)
	}

	server := server.CreateServer(conf)
	err = server.Run()
	if err != nil {
		log.Fatal().Err(err)
	}
}
