package locked

import (
	"context"

	"github.com/0x0FACED/locked/internal/app/services"
)

type app struct {
	file services.SecretService
}

func New() *app {
	return &app{}
}

func StartCLI(ctx context.Context) error {
	panic("impl me")
}

func StartWeb(ctx context.Context) error {
	panic("impl me")
}
