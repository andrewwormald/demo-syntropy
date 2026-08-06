package main

import (
	"fmt"

	"connect4-go/board"
)

func main() {
	fmt.Println("Connect 4")
	b := board.New()
	fmt.Print(board.Render(b))
}
