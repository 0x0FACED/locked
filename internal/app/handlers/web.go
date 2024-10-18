package handlers

import "github.com/0x0FACED/locked/internal/app/services"

type WebHandler struct {
	secretSrv services.SecretService
}

func NewWeb(srv services.SecretService) WebHandler {
	return WebHandler{
		secretSrv: srv,
	}
}
