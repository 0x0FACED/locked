package locked

import (
	"context"
	"flag"
	"fmt"
)

// Здесь мы обрабатываем флаги и вызываем метод сервиса
func (a *cliApp) add(ctx context.Context, args []string) {
	// Создание и парсинг флагов для команды add
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)

	// Надо ыб добавить флаг каокй-то, который будет указывать на фоновый режим. Если фоновый режим,
	// то в горутине запускаем процесс, если же не фоновый, то блкокируем основной ввод и ждем.

	// флаги
	name := addCmd.String("n", "", "Name of the secret")
	desc := addCmd.String("d", "", "Description of the secret")
	payload := addCmd.String("s", "", "Secret payload or file path")

	// парсинг аргументов
	if err := addCmd.Parse(args); err != nil {
		fmt.Println("Error parsing flags:", err)
		a.errCh <- err
	}

	// проверка, что все обязательные флаги заполнены
	if *name == "" || *desc == "" || *payload == "" {
		fmt.Println("Usage: add -n <name> -d <description> -s <payload>")
		return
	}

	// создание параметров команды
	//params := models.AddSecretCmdParams{
	//	Name:        *name,
	//	Description: *desc,
	//	Payload:     *payload,
	//}

	// вызов метода для добавления секрета
	//a.secretService.Add(ctx, params)
}

func (a *cliApp) open(ctx context.Context, filename string) {
	a.secretService.Open(ctx, filename)
}
