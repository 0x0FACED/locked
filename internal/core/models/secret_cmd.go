package models

import "github.com/0x0FACED/locked/internal/core/models/types"

// Эта структура нужна для отображения секрета в консоли/web/desktop UI

type SecretUI struct {
	ID          types.SecretID
	Name        types.SecretName
	Description types.SecretDescription
	Type        types.SecretType
	CreatedAt   types.SecretCreatedAt
	Size        types.SecretSize
	// Payload как таковая не будет отображаться, если секрет - файл или длинный текст.
	Payload types.SecretPayload
}
