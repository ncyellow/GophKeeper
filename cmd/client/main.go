// Package main содержит запуск клиента.
package main

import (
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/client/config"
	"github.com/ncyellow/GophKeeper/internal/client/console"
)

func main() {
	fmt.Println("Client run")
	terminal := console.Console{
		Conf: config.ParseConfig(),
	}
	terminal.Run()
}
