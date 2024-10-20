package handlers

import (
	"net/http"

	"github.com/0x0FACED/locked/internal/app/services"
)

type WebHandler struct {
	server    http.Server
	secretSrv services.SecretService
}

func NewWeb(srv services.SecretService) *WebHandler {
	return &WebHandler{
		server:    http.Server{}, // временно
		secretSrv: srv,
	}
}
