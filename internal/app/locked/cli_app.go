package locked

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/0x0FACED/locked/internal/app/services"
	"github.com/0x0FACED/locked/internal/core/models"
	"github.com/chzyer/readline"
	"golang.org/x/term"
)

type cliApp struct {
	secretService services.SecretService
	currentFile   *os.File
}

// Реализация автодополнения команд
type autoCompleter struct{}

func (a autoCompleter) Do(line []rune, pos int) (newLine [][]rune, offset int) {
	commands := []string{"add", "open", "close", "exit"}
	for _, cmd := range commands {
		if strings.HasPrefix(string(line), cmd) {
			newLine = append(newLine, []rune(cmd))
		}
	}
	return newLine, len(line)
}

func NewCLIApp(resCh chan []byte, errCh chan error, done chan struct{}) *cliApp {
	secretService := services.New(resCh, errCh, done)
	return &cliApp{
		secretService: secretService,
	}
}

func (a *cliApp) StartCLI(ctx context.Context) error {
	// Инициализация autocomplete с командами
	commands := []string{"add", "open", "clear", "close", "del", "exit"}

	// Создание обработчика для автодополнения
	completer := readline.NewPrefixCompleter(
		readline.PcItem("add"),
		readline.PcItem("open"),
		readline.PcItem("clear"),
		readline.PcItem("close"),
		readline.PcItem("del"),
		readline.PcItem("exit"),
	)

	// Инициализируем readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "locked ~# ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer, // Поддержка автозаполнения
	})
	if err != nil {
		return err
	}

	defer rl.Close()

	// Каналы для обработки ввода и ошибок
	inputCh := make(chan string)
	errCh := make(chan error)

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
			rl.SetPrompt(fmt.Sprintf("locked/%s ~# ", a.currentFile.Name()))
		} else {
			rl.SetPrompt("locked ~# ")
		}

		go readCmd(rl, inputCh, errCh)

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
					// сет флагов для команды add
					addCmd := flag.NewFlagSet("add", flag.ContinueOnError)

					// флаги
					// TODO: как-то это все структурировать в одном месте, а не при каждом вызове функции
					name := addCmd.String("n", "", "Name of the secret")
					desc := addCmd.String("d", "", "Description of the secret")
					payload := addCmd.String("s", "", "Secret payload or file path")

					// парсинг
					if err := addCmd.Parse(words[1:]); err != nil {
						fmt.Println("Error parsing flags:", err)
						return nil
					}

					// проверяем, что обязательные флаги есть (они все обязательные лол)
					if *name == "" || *desc == "" || *payload == "" {
						fmt.Println("Usage: add -n <name> -d <description> -s <payload>")
						return nil
					}

					// для удобства юзаем структуру параметров
					params := models.AddSecretCmdParams{
						Name:        *name,
						Description: *desc,
						Payload:     *payload, // это может быть путь к файлу или текстовый payload
					}

					// Вызываем метод для добавления секрета
					a.secretService.Add(ctx, params)
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
				err := a.currentFile.Close()
				if err != nil {
					fmt.Println("Error closing the file:", err)
				}
				a.currentFile = nil
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

func readCmd(rl *readline.Instance, inputCh chan string, errCh chan error) {
	line, err := rl.Readline()
	if err != nil {
		if len(line) != 0 {
			errCh <- err
		}
	}
	inputCh <- line
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
