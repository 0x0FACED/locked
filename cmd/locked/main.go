package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0x0FACED/locked/internal/app/locked"
	"github.com/0x0FACED/locked/internal/core/models"
)

const (
	secretsDir     = "secrets"
	masterHashFile = "master_hash"
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

	isFirst := checkFirstRun()

	command := os.Args[1]
	resCh := make(chan models.Result, 10) // временно 10, потом думаю через конфиг передавать
	errCh := make(chan error, 5)          // временно 5, потом думаю через конфиг передавать
	done := make(chan struct{}, 2)        // временно 2, потом думаю через конфиг передавать

	// Пока что так сделал, но это не совсем гуд, как мне кажется
	if command == cli {
		app := locked.NewCLIApp(resCh, errCh, done)
		app.StartCLI(ctx, isFirst)
	} else if command == web {
		app := locked.NewWebApp(resCh, errCh, done)
		app.StartWeb(ctx, isFirst)
	} else {
		fmt.Println("Invalid command. Use 'locked cli' or 'locked web'.")
		os.Exit(1)
	}
}

func checkFirstRun() bool {
	if _, err := os.Stat(filepath.Join(secretsDir, masterHashFile)); os.IsNotExist(err) {
		return true // файл не найден - это первый запуск
	}
	return false
}
