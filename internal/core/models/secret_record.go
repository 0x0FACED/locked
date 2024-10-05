package models

import "github.com/0x0FACED/locked/internal/core/models/types"

// Запись секрета в файле
type SecretRecord struct {
	Offset      types.SecretOffset      // Смещение в файле
	ID          types.SecretID          // Уникальный идентификатор секрета
	Name        types.SecretName        // Название секрета
	Type        types.SecretType        // Тип секрета (текст, медиа и т.д.)
	CreatedAt   types.SecretCreatedAt   // Дата создания
	Size        types.SecretSize        // Размер секрета
	Description types.SecretDescription // Описание секрета
	Payload     types.SecretPayload     // Зашифрованный секрет
}

// Порядок записи данных полей в файл еще не уточнен
