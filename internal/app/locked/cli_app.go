package locked

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/0x0FACED/locked/internal/app/services"
	"github.com/0x0FACED/locked/internal/core/models"
	"github.com/chzyer/readline"
)

const VESRION = uint8(1) // пишем в header в первую очередь версию
const EXTENSION = "lkd"

var BASE_PKG = "secrets"

type cliApp struct {
	currentFile string
	rl          *readline.Instance

	secretService services.SecretService

	//wp     *worker.WorkerPool
	//taskCh chan worker.Task

	resCh chan models.Result
	errCh chan error
	done  chan struct{}

	updPromptCh chan string
}

func NewCLIApp(resCh chan models.Result, errCh chan error, done chan struct{}) *cliApp {
	secretService := services.New(resCh, errCh, done)
	completer := completer()

	//taskCh := make(chan worker.Task, 10) // 10 задач пока что пускай

	//wp := worker.New(secretService, taskCh, errCh) // создаем воркер пул

	updPromptCh := make(chan string)
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
		rl:            rl, // rl для работы со строкой ввода
		secretService: secretService,
		//wp:     wp,     // вп будет чекать taskCh и обрабатывать задачи
		//taskCh: taskCh, // сюда сливать будем все задачи
		resCh: resCh, // мб удалю
		errCh: errCh, // мб удалю
		done:  done,  // мб удалю

		updPromptCh: updPromptCh,
	}
}

// автокомплит при вводе части команды и нажатии на tab
func completer() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("new"),
		readline.PcItem("ls"),
		readline.PcItem("add"),
		readline.PcItem("open"),
		readline.PcItem("clear"),
		readline.PcItem("close"),
		readline.PcItem("del"),
		readline.PcItem("exit"),
	)
}

// основной метод для запуска всего приложения
func (a *cliApp) StartCLI(ctx context.Context, isFirstRun bool) {
	if isFirstRun {
		fmt.Println("~ welcome, samurai ~")
		fmt.Println("~ Before you start, you need to create a password ~")
		fmt.Println("~ !IMPORTANT! This password is your main key to all your secrets ~")
		err := initApp()
		if err != nil {
			fmt.Println("~ error happened:", err)
		}
	} else {
		verify()
	}

	//a.wp.Start(ctx)

	fmt.Println("~ ~ ~ welcome back, samurai ~ ~ ~")

	// запускаем отдельную горутину для прослушивая результатов выполнения команд
	// open, close, clear etc
	go a.listen()

	a.run(ctx)
}

func (a *cliApp) run(ctx context.Context) {
	inputCh := make(chan string)
	errCh := make(chan error)

	for {
		if a.currentFile != "" {
			a.rl.SetPrompt(fmt.Sprintf("locked/%s ~# ", a.currentFile))
		} else {
			a.rl.SetPrompt("locked ~# ")
		}

		// горутина для ввода данных
		go readCmd(a.rl, inputCh, errCh)

		select {
		case input := <-inputCh: // обработка команды
			a.handleCommand(ctx, input) // обрабатываем команду
		case err := <-errCh: // ошибка при вводе команды
			fmt.Println("~ Error reading input:", err)
		case newPrompt := <-a.updPromptCh: // обновление промпта
			a.rl.SetPrompt(newPrompt)
			a.rl.Refresh() // перерисовка строки
		}

	}
}

func isFileExists(filename string) bool {
	fullName := filename + "." + EXTENSION
	if _, err := os.Stat(filepath.Join(secretsDir, fullName)); err == nil {
		fmt.Println("~ File with this name already exists")
		return true
	} else if !os.IsNotExist(err) {
		// В случае ошибки, отличной от "файл не существует"
		fmt.Println("~ Something went wrong with error:", err)
		return false
	}

	return false
}

func createSecretFile(filename string) error {
	// Если файл не существует, создаем его
	fullName := filename + ".lkd"
	file, err := os.Create(filepath.Join(secretsDir, fullName))
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	fmt.Printf("~ File %s created successfully in %s\n", fullName, secretsDir)
	return nil
}

func (a *cliApp) handleCommand(ctx context.Context, input string) {
	command := strings.TrimSpace(input)
	words := strings.Split(command, " ")
	switch words[0] {
	case "new":
		if !isFileExists(words[1]) {
			err := createSecretFile(words[1])
			if err != nil {
				fmt.Println("~ Error creating file with error:", err)
				return
			}

		}
	case "add": // добавление секрета
		if a.checkFileStatus() != nil {
			// никакой файл не открыт, добавлять некуда!
			fmt.Println("~ You need to open any of your secret files or create one to keep a secret.")
			fmt.Println("~ To open the file, type the following command: open filename.lkd")
		} else {
			a.add(ctx, words)
			/*task := worker.Task{
				Command: words[0], // open
				Args:    words[1],
			}*/

			//a.taskCh <- task
		}

	case "open":
		if len(words) != 2 {
			fmt.Println("~ To open the file, type the following command: open filename.lkd")
			break
		}

		a.open(ctx, words[1])

		/*task := worker.Task{
			Command: words[0], // open
			Args:    words[1],
		}*/

		//a.taskCh <- task
		/*
			f, err := os.Open(words[1]) // заглушка
			if err != nil {
				fmt.Println("err:", err)
			}
			a.currentFile = f
		*/
	case "clear": // очистка всего файла с секретами

	case "close": // закрыть файл, но не приложение
		/*
			err := a.currentFile.Close()
			if err != nil {
				fmt.Println("Error closing the file:", err)
			}
		*/
		a.currentFile = "" // xD закрыли)))
	case "del": // удалить секрет из файла
	case "exit": // выход из приложения
		fmt.Println("~ ~ ~ Exiting the application. Goodbye! ~ ~ ~")
		/*
			err := a.currentFile.Close()
			if err != nil {
				fmt.Println("Error closing the file:", err)
			}
		*/
		os.Exit(0)
	}
}

func (a *cliApp) listen() {
	for {
		select {
		case result := <-a.resCh:
			switch result.Command {
			case "add":
				// ..
			case "open":
				fmt.Printf("~ File %s opened.\n", string(result.Data))
				a.currentFile = string(result.Data)

				a.updPromptCh <- fmt.Sprintf("locked/%s ~# ", a.currentFile)
			}
		case err := <-a.errCh:
			fmt.Println("~ Error:", err) // вывод ошибки
		case <-a.done:
			fmt.Println("~ Task completed!") // сигнал о завершении
		}
	}
}

func readCmd(rl *readline.Instance, inputCh chan string, errCh chan error) {
	line, err := rl.Readline()
	if err != nil {
		if err == readline.ErrInterrupt { // ловим Ctrl+C
			fmt.Println("~ Exiting...")
			os.Exit(0)
		}
		if len(line) != 0 {
			errCh <- err
		}
	}
	inputCh <- line
}

func (a *cliApp) checkFileStatus() error {
	if a.currentFile == "" {
		return os.ErrClosed // немного неправильно, ибо в доке написано "file ALREADY closed", но пока что так
	}
	return nil
}
