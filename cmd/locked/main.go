package main

import (
	"context"
	"fmt"
	"os"

	"github.com/0x0FACED/locked/internal/app/locked"
)

var (
	cli = "cli"
	web = "web"
)

func main() {
	ctx := context.Background()

	if len(os.Args) != 2 {
		fmt.Println("Run with command: locked [cli|web]")
		fmt.Println("$ locked cli // starts cli ui")
		fmt.Println("$ locked web // starts web ui")
		os.Exit(1)
	}

	command := os.Args[1]
	resCh := make(chan []byte, 10) // временно 10, потом думаю через конфиг передавать
	errCh := make(chan error, 5)   // временно 5, потом думаю через конфиг передавать
	done := make(chan struct{}, 2) // временно 2, потом думаю через конфиг передавать

	// Пока что так сделал, но это не совсем гуд, как мне кажется
	if command == cli {
		app := locked.NewCLIApp(resCh, errCh, done)
		app.StartCLI(ctx)
	} else if command == web {
		app := locked.NewWebApp(resCh, errCh, done)
		app.StartWeb(ctx)
	} else {
		fmt.Println("Invalid command. Use 'locked cli' or 'locked web'.")
		os.Exit(1)
	}
}
