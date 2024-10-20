package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/0x0FACED/locked/internal/app/locked"
)

var (
	cli = "cli"
	web = "web"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	if len(os.Args) != 2 {
		fmt.Println("Run with command: locked [cli|web]")
		fmt.Println("$ locked cli // starts cli ui")
		fmt.Println("$ locked web // starts web ui")
		os.Exit(1)
	}

	command := os.Args[1]

	// Пока что так сделал, но это не совсем гуд, как мне кажется
	if command == cli {
		app := locked.NewCLIApp()
		app.StartCLI(ctx)
	} else if command == web {
		app := locked.NewWebApp()
		app.StartWeb(ctx)
	} else {
		fmt.Println("Invalid command. Use 'locked cli' or 'locked web'.")
		os.Exit(1)
	}
}
