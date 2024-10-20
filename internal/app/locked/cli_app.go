package locked

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/0x0FACED/locked/internal/app/handlers"
	"github.com/0x0FACED/locked/internal/app/services"
	"golang.org/x/term"
)

type cliApp struct {
	cliHandler  *handlers.CLIHandler
	currentFile *os.File
}

func NewCLIApp() *cliApp {
	secretService := services.New()
	return &cliApp{
		cliHandler: handlers.NewCLI(secretService),
	}
}

func (a *cliApp) StartCLI(ctx context.Context) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		err := requestPassword()
		if err != nil {
			fmt.Println("Incorrect password, try again")
			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("Successfully logged in!")
			break
		}
	}
	fmt.Println("Welcome to the CLI Application! Type 'exit' to quit.")

	for {
		if a.currentFile != nil {
			fmt.Printf("locked/%s ~# ", a.currentFile.Name())
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
		case <-ctx.Done(): // graceful shutdown
			return nil
		case input := <-inputCh: // обработка команды
			command := strings.TrimSpace(input)
			words := strings.Split(command, " ")
			switch words[0] {
			case "add": // добавление секрета
				if a.checkFileStatus() != nil {
					// никакой файл не открыт, добавлять некуда!
					fmt.Println("You need to open any of your secret files or create one to keep a secret.")
					fmt.Println("To open the file, type the following command: open filename.lkd")
				} else {
					fmt.Println("You entered:", words[1])
				}

			case "open":
				if len(words) != 2 {
					fmt.Println("To open the file, type the following command: open filename.lkd")
					continue
				}
				f, err := os.Open(words[1]) // заглушка
				if err != nil {
					fmt.Println("err:", err)
				}
				a.currentFile = f
			case "clear": // очистка всего файла с секретами

			case "close": // закрыть файл, но не приложение

			case "del": // удалить секрет из файла
			case "exit": // выход из приложения
				fmt.Println("Exiting the application. Goodbye!")
				a.currentFile.Close()
				return nil
			}

		case err := <-errCh: // ошибка при вводе команды
			fmt.Println("Error reading input:", err)
		}
	}
}

// запрашиваем ввод пароля
func requestPassword() error {
	fmt.Print("Enter master password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Println() // Для переноса строки после ввода пароля
	password := string(bytePassword)
	if password == "admin" { // ВРЕМЕННО заглушка
		return nil
	}
	return errors.New("incorrect") // TODO: сделать адекватно
}

func (a *cliApp) checkFileStatus() error {
	if a.currentFile == nil {
		return os.ErrClosed // немного неправильно, ибо в доке написано "file ALREADY closed", но пока что так
	}
	return nil
}
