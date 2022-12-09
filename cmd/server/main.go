package main

import (
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/server/config"
	"github.com/ncyellow/GophKeeper/internal/server/httpserver"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Println("Server start")

	server := httpserver.HTTPServer{
		Conf: config.ParseConfig(),
	}
	server.Run()
}
