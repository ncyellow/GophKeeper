// Package main содержит запуск сервера.
package main

import (
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/server"
	"github.com/ncyellow/GophKeeper/internal/server/config"
)

func main() {
	fmt.Println("Server start")
	server := server.CreateServer(config.ParseConfig())
	server.Run()
}
