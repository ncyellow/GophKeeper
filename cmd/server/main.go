package main

import (
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/server"
	"github.com/ncyellow/GophKeeper/internal/server/config"
)

func main() {
	fmt.Printf("Build version: %s\n", config.BuildVersion)
	fmt.Printf("Build date: %s\n", config.BuildDate)
	fmt.Println("Server start")

	server := server.CreateServer(config.ParseConfig())
	server.Run()
}
