package locked

import (
	"context"
	"os"

	"github.com/0x0FACED/locked/internal/app/services"
	"github.com/0x0FACED/locked/internal/core/models"
)

type webApp struct {
	secretService services.SecretService
	currentFile   *os.File
}

func NewWebApp(resCh chan models.Result, errCh chan error, done chan struct{}) *webApp {
	secretService := services.New(resCh, errCh, done)
	return &webApp{
		secretService: secretService,
	}
}

func (a *webApp) StartWeb(ctx context.Context, isFirstRun bool) error {
	panic("impl me")
}
