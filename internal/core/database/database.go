package database

import (
	"context"
	"os"

	"github.com/0x0FACED/locked/internal/core/models"
)

type helper interface {
	Offset() int64
}

// Интерфейс для работы с файлами (чтение файла целиком, запись в файл,
// очистка файла целиком, удаление какой-то записи
type Database interface {
	Open(ctx context.Context, filename string) (*os.File, error)
	// Read читает определенный секрет (хз, зачем, но мб пригодится)
	Read()
	// ReadFull читает фулл содержимое файла
	ReadFull()
	// Write записывает секрет в конец файла
	Write(ctx context.Context, secret models.SecretRecord)
	// Flush очищает файл полностью, но файл не удаляет
	Flush()
	// Delete удаляет определенный секрет
	Delete()
	// DeleteN удаляет некоторое количество секретов (перечисленные через ,)
	DeleteN()

	// какие-то методы в теории
	helper
}

type fileDatabase struct {
	currentFile *os.File
	currPos     int64
}

func (d *fileDatabase) Offset() int64 {
	return d.currPos
}
