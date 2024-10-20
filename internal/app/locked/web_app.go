package locked

import (
	"context"
	"os"

	"github.com/0x0FACED/locked/internal/app/handlers"
	"github.com/0x0FACED/locked/internal/app/services"
)

type webApp struct {
	webHandler  *handlers.WebHandler
	currentFile *os.File
}

func NewWebApp() *webApp {
	secretService := services.New()
	return &webApp{
		webHandler: handlers.NewWeb(secretService),
	}
}

func (a *webApp) StartWeb(ctx context.Context) error {
	panic("impl me")
}
