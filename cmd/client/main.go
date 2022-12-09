package main

import (
	"fmt"

	"github.com/ncyellow/GophKeeper/internal/client/console"
)

func main() {
	fmt.Println("Client run")
	terminal := console.Console{}
	terminal.Run()
}
