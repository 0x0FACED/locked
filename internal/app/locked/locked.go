package locked

import (
	"github.com/0x0FACED/locked/internal/app/services"
)

type app struct {
	file services.SecretService
}

func New() *app {
	return &app{}
}

func Start() error {
	panic("impl me")
}
