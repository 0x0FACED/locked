package locked

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/0x0FACED/locked/internal/app/services"
	"github.com/chzyer/readline"
	"golang.org/x/term"
)

var BASE_PKG = "secrets/"

type cliApp struct {
	secretService services.SecretService
	currentFile   *os.File
	rl            *readline.Instance
	resCh         chan []byte
	errCh         chan error
	done          chan struct{}
}

func NewCLIApp(resCh chan []byte, errCh chan error, done chan struct{}) *cliApp {
	secretService := services.New(resCh, errCh, done)
	completer := completer()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "locked ~# ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer, // автодополнение
	})

	if err != nil {
		return nil
	}

	return &cliApp{
		secretService: secretService,
		rl:            rl,
		resCh:         resCh,
		errCh:         errCh,
		done:          done,
	}
}

func completer() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("add"),
		readline.PcItem("open"),
		readline.PcItem("clear"),
		readline.PcItem("close"),
		readline.PcItem("del"),
		readline.PcItem("exit"),
	)
}

func (a *cliApp) StartCLI(ctx context.Context) {

	// отдельная функция, просто бесконечный цикл для логина (пароль)
	verify()

	fmt.Println("~ ~ ~ welcome back, samurai ~ ~ ~")

	// запускаем отдельную горутину для прослушивая результатов выполнения команд
	// open, close, clear etc
	go a.listen()

	a.run(ctx)
}

func verify() {
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
}

func (a *cliApp) run(ctx context.Context) {
	inputCh := make(chan string)
	errCh := make(chan error)

	for {
		if a.currentFile != nil {
			a.rl.SetPrompt(fmt.Sprintf("locked/%s ~# ", a.currentFile.Name()))
		} else {
			a.rl.SetPrompt("locked ~# ")
		}

		// горутина для ввода данных
		go readCmd(a.rl, inputCh, errCh)

		select {
		case input := <-inputCh: // обработка команды
			a.handleCommand(ctx, input)
		case err := <-errCh: // ошибка при вводе команды
			fmt.Println("Error reading input:", err)
		}
	}
}

func (a *cliApp) handleCommand(ctx context.Context, input string) {
	command := strings.TrimSpace(input)
	words := strings.Split(command, " ")
	switch words[0] {
	case "add": // добавление секрета
		if a.checkFileStatus() != nil {
			// никакой файл не открыт, добавлять некуда!
			fmt.Println("You need to open any of your secret files or create one to keep a secret.")
			fmt.Println("To open the file, type the following command: open filename.lkd")
		} else {
			go a.add(ctx, words) // щас в горутине отдельно, чтобы заблочить управление юзеру
		}

	case "open":
		if len(words) != 2 {
			fmt.Println("To open the file, type the following command: open filename.lkd")
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
	}
}

func (a *cliApp) listen() {
	for {
		select {
		case result := <-a.resCh:
			fmt.Println("Result:", string(result)) // вывод результата
		case err := <-a.errCh:
			fmt.Println("Error:", err) // вывод ошибки
		case <-a.done:
			fmt.Println("Task completed!") // сигнал о завершении
		}
	}
}

func readCmd(rl *readline.Instance, inputCh chan string, errCh chan error) {
	line, err := rl.Readline()
	if err != nil {
		if err == readline.ErrInterrupt { // ловим Ctrl+C
			fmt.Println("Exiting...")
			os.Exit(0)
		}
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
