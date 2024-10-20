package handlers

import (
	"context"

	"github.com/0x0FACED/locked/internal/app/services"
)

type CLIHandler struct {
	secretService services.SecretService
}

func NewCLI(srv services.SecretService) *CLIHandler {
	return &CLIHandler{
		secretService: srv,
	}
}
