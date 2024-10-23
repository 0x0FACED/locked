package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/0x0FACED/locked/internal/app/services"
	"github.com/0x0FACED/locked/internal/core/models"
)

// Вообще суть в том, что все таски будут идти в один канал, по сути формируя очередь
// Эти таски будут в отдельной горутине читаться и на каждую выделяться отдельная горутина (воркер)
// Количество воркеров будет думаю в конфиге прописываться

// структура для задачи
type Task struct {
	// команда (add, open, close etc)
	Command string
	// аргументы (включая и флаги	)
	Args string
}

type WorkerPool struct {
	secretService services.SecretService
	wg            sync.WaitGroup

	MaxWorkers int
	TaskCh     chan Task

	errCh chan error
}

func New(service services.SecretService, taskQueue chan Task, errCh chan error) *WorkerPool {
	return &WorkerPool{
		secretService: service,
		MaxWorkers:    10, // Временно
		TaskCh:        taskQueue,
		// чтобы вернуть ошибку сразу, если неправильная команда введена
		errCh: errCh,
	}
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.MaxWorkers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.work(ctx)
		}()
	}
}

func (p *WorkerPool) work(ctx context.Context) {
	for task := range p.TaskCh {
		switch task.Command {
		case "add":
			p.secretService.Add(ctx, models.AddSecretCmdParams{})

		case "open":
			// пока что без валидации
			p.secretService.Open(ctx, task.Args)
			// Добавить остальные команды

		default:
			p.errCh <- errors.New("test error: unknown command")
		}
	}
}
