package database

import (
	"context"
	"errors"
	"os"
)

type helper interface {
	Offset() uint64
	Count() uint8
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
	Write(ctx context.Context, secret []byte) error
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
	currPos     uint64
	secretCnt   uint8
}

func NewFileDatabase() Database {
	return &fileDatabase{}
}

func (d *fileDatabase) Open(ctx context.Context, filename string) (*os.File, error) {
	panic("not implemented") // TODO: Implement
}

func (d *fileDatabase) Read() {
	panic("not implemented") // TODO: Implement
}

func (f *fileDatabase) ReadFull() {
	panic("not implemented") // TODO: Implement
}

func (d *fileDatabase) Write(ctx context.Context, secret []byte) error {
	n, err := d.currentFile.Write(secret)
	if err != nil {
		return err
	}
	if n != len(secret) {
		return errors.New("~ Written data size not equal to secret size")
	}

	return nil
}

func (d *fileDatabase) Flush() {
	panic("not implemented") // TODO: Implement
}

func (d *fileDatabase) Delete() {
	panic("not implemented") // TODO: Implement
}

func (d *fileDatabase) DeleteN() {
	panic("not implemented") // TODO: Implement
}

func (d *fileDatabase) Offset() uint64 {
	return d.currPos
}

func (d *fileDatabase) Count() uint8 {
	return d.secretCnt
}
