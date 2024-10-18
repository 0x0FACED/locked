package adapters

import "github.com/0x0FACED/locked/internal/app/services"

type WebAdapter struct {
	secretSrv services.SecretService
}

func NewWeb(srv services.SecretService) WebAdapter {
	return WebAdapter{
		secretSrv: srv,
	}
}
