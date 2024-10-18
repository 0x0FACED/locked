package handlers

import "github.com/0x0FACED/locked/internal/app/services"

type CLIHandler struct {
	secretSrv services.SecretService
}

func NewCLI(srv services.SecretService) CLIHandler {
	return CLIHandler{
		secretSrv: srv,
	}
}
