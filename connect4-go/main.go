package main

import "fmt"

const (
	rows = 6
	cols = 7
)

func printBoard() {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			fmt.Print("| . ")
		}
		fmt.Println("|")
	}
	for c := 0; c < cols; c++ {
		fmt.Print("----")
	}
	fmt.Println("-")
}

func main() {
	fmt.Println("Connect 4")
	printBoard()
}
