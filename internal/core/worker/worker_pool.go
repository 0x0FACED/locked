package worker

import (
	"context"
	"sync"

	"github.com/0x0FACED/locked/internal/app/services"
)

// Вообще суть в том, что все таски будут идти в один канал, по сути формируя очередь
// Эти таски будут в отдельной горутине читаться и на каждую выделяться отдельная горутина (воркер)
// Количество воркеров будет думаю в конфиге прописываться

// структура для задачи
type Task struct {
	// команда (add, open, close etc)
	Command string
	// аргументы (включая и флаги	)
	Args []string
}

type WorkerPool struct {
	secretService services.SecretService
	wg            sync.WaitGroup

	MaxWorkers int
	TaskCh     chan Task
}

func New(service services.SecretService, taskQueue chan Task) *WorkerPool {
	return &WorkerPool{
		secretService: service,
		MaxWorkers:    10, // Временно
		TaskCh:        taskQueue,
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
			p.secretService.Add(ctx, task.Args)

			// Добавить остальные команды
		}
	}
}
