package main

import (
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/httpserver"
)

func main() {
	fmt.Printf("Build version: %s\n", config.BuildVersion)
	fmt.Printf("Build date: %s\n", config.BuildDate)
	fmt.Println("Server start")

	server := httpserver.HTTPServer{
		Conf: config.ParseConfig(),
	}
	server.Run()
}
