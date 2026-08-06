package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"connect4-go/game"
)

func main() {
	fmt.Println("Connect 4")

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "connect4:", err)
		os.Exit(1)
	}
	defer term.Restore(fd, oldState)

	// Raw mode disables normal signal delivery for interrupts received while
	// the terminal has focus; catch them explicitly so the terminal is
	// restored to cooked mode before the process exits.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		term.Restore(fd, oldState)
		os.Exit(1)
	}()

	if err := game.Run(os.Stdin, os.Stdout); err != nil {
		term.Restore(fd, oldState)
		fmt.Fprintln(os.Stderr, "connect4:", err)
		os.Exit(1)
	}
}
