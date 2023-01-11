// Package main содержит запуск сервера.
package main

import (
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/server"
	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/rs/zerolog/log"
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
