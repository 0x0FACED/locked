package models

// Запись секрета в файле
// 8 + 1 + 64 + 1 + 8 + 8 + 128 + n = 200 + 18
type SecretRecord struct {
	Offset      uint64    // Смещение в файле
	Name        [64]byte  // Название секрета
	Type        uint8     // Тип секрета (текст, медиа и т.д.)
	CreatedAt   uint64    // Дата создания
	Size        uint64    // Размер секрета
	Description [128]byte // Описание секрета
	Payload     []byte    // Зашифрованный секрет
}

// Порядок записи данных полей в файл еще не уточнен
