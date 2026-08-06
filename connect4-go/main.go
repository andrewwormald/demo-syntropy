package main

import (
	"fmt"
	"os"

	"connect4-go/game"
)

func main() {
	fmt.Println("Connect 4")
	if err := game.Run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "connect4:", err)
		os.Exit(1)
	}
}
