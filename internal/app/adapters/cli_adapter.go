package adapters

import "github.com/0x0FACED/locked/internal/app/services"

type CLIAdapter struct {
	secretSrv services.SecretService
}

func NewCLI(srv services.SecretService) CLIAdapter {
	return CLIAdapter{
		secretSrv: srv,
	}
}
