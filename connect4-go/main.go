package main

import (
	"fmt"
	"strings"
)

const (
	rows = 6
	cols = 7
)

func renderBoard() string {
	var b strings.Builder
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			b.WriteString("| . ")
		}
		b.WriteString("|\n")
	}
	for c := 0; c < cols; c++ {
		b.WriteString("----")
	}
	b.WriteString("-\n")
	return b.String()
}

func main() {
	fmt.Println("Connect 4")
	fmt.Print(renderBoard())
}
