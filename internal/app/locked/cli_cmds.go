package locked

import (
	"context"
	"errors"
	"flag"

	"github.com/0x0FACED/locked/internal/core/models"
)

const (
	CLOSE = "close"
	ADD   = "add"
	OPEN  = "open"
	CLEAR = "clear"
	EXIT  = "exit"
	NEW   = "new"
	DEL   = "del"
)

// Здесь мы обрабатываем флаги и вызываем метод сервиса
func (a *cliApp) add(ctx context.Context, args []string) {
	// Создание и парсинг флагов для команды add
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)

	// Надо ыб добавить флаг каокй-то, который будет указывать на фоновый режим. Если фоновый режим,
	// то в горутине запускаем процесс, если же не фоновый, то блкокируем основной ввод и ждем.

	// флаги
	name := flag.String("n", "", "Name of the secret")
	desc := flag.String("d", "", "Description of the secret")
	secretText := flag.String("s", "", "Secret text data")
	secretFile := flag.String("p", "", "Path to file containing secret")

	// парсинг аргументов
	if err := addCmd.Parse(args); err != nil {
		a.errCh <- err
		return
	}

	if *secretText != "" && *secretFile != "" {
		a.errCh <- errors.New("~ Error: Please specify only one of -s (secret text) or -p (secret file path)")
		return
	}

	// НЕ ЗДЕСЬ
	/*
		// проверка, что все обязательные флаги заполнены
		if *name == "" || *desc == "" || *payload == "" {
			fmt.Println("Usage: add -n <name> -d <description> -s <payload>")
			return
		}
	*/

	// инициализация параметров для передачи в сервис
	var params models.AddSecretCmdParams
	params.Name = name
	params.Description = desc

	// Проверка на тип секрета (текст или файл) и установка параметров
	if *secretText != "" {
		params.Payload = secretText
		params.IsFile = false
	} else if *secretFile != "" {
		params.Payload = secretFile
		params.IsFile = true
	} else {
		a.errCh <- errors.New("~ Error: Either -s (secret text) or -p (secret file path) must be provided")
		return
	}

	// вызов метода для добавления секрета
	a.secretService.Add(ctx, params)
}

func (a *cliApp) open(ctx context.Context, filename string) {
	a.secretService.Open(ctx, filename)
}

func (a *cliApp) createSecretFile(ctx context.Context, filename string) {
	a.secretService.CreateSecretFile(ctx, filename)
}
