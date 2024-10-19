package locked

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/0x0FACED/locked/internal/app/services"
)

type app struct {
	file        services.SecretService
	currentFile *os.File
}

func new() *app {
	return &app{
		file: services.New(),
	}
}

func StartCLI(ctx context.Context) error {
	reader := bufio.NewReader(os.Stdin)
	app := new()
	fmt.Println("Welcome to the CLI Application! Type 'exit' to quit.")

	for {
		if app.currentFile != nil {
			fmt.Printf("locked/%s ~# ", app.currentFile.Name())
		} else {
			fmt.Printf("locked ~# ")
		}

		inputCh := make(chan string, 1)
		errCh := make(chan error, 1)

		// читаем из консоли в отдельной горутине
		go func() {
			input, err := reader.ReadString('\n')
			if err != nil {
				errCh <- err
			} else {
				inputCh <- input
			}
		}()

		select {
		case <-ctx.Done():
			return nil
		case input := <-inputCh:
			command := strings.TrimSpace(input)
			words := strings.Split(command, " ")
			switch words[0] {
			case "add": // добавление секрета

			case "open":
				f, err := os.Open("README.md") // заглушка
				if err != nil {
					fmt.Println("err:", err)
				}
				app.currentFile = f
			case "clear": // очистка всего файла с секретами

			case "close": // закрыть файл, но не приложение

			case "del": // удалить секрет из файла
			case "exit": // выход из приложения
			}
			if command == "exit" {
				fmt.Println("Exiting the application. Goodbye!")
				return nil
			}

		case err := <-errCh:
			fmt.Println("Error reading input:", err)
		}
	}
}

func StartWeb(ctx context.Context) error {
	panic("impl me")
}
