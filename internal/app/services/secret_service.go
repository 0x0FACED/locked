package services

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/0x0FACED/locked/internal/core/database"
	"github.com/0x0FACED/locked/internal/core/encryption"
	"github.com/0x0FACED/locked/internal/core/models"
)

const (
	EXTENSION   = "lkd"
	SECRETS_DIR = "secrets"
)

type SecretService interface {
	Add(ctx context.Context, secret models.AddSecretCmdParams)
	Open(ctx context.Context, filename string)
	CreateSecretFile(ctx context.Context, filename string)
}

type secretService struct {
	currentFile *os.File

	zip   zip.Compressor
	unzip zip.Decompressor

	enc encryption.Encryptor
	dec encryption.Decryptor

	db database.Database

	resCh chan models.Result
	errCh chan error
	done  chan struct{}
}

func New(resCh chan models.Result, errCh chan error, done chan struct{}) SecretService {
	return &secretService{
		resCh: resCh,
		errCh: errCh,
		done:  done,
	}
}

func (s *secretService) Open(ctx context.Context, filename string) {
	f, err := os.Open(filename) // только для чтения открываем пока что
	if err != nil {
		s.errCh <- err
	}
	s.currentFile = f

	res := models.Result{
		Command: "open",
		Data:    []byte(f.Name()),
	}

	s.resCh <- res
}

func (s *secretService) CreateSecretFile(ctx context.Context, filename string) {
	if isFileExists(filename) {
		s.errCh <- errors.New("file already exists")
		return
	}
	// Если файл не существует, создаем его
	fullName := filename + "." + EXTENSION
	file, err := os.Create(filepath.Join(SECRETS_DIR, fullName))
	if err != nil {
		s.errCh <- err
		return
	}

	defer file.Close()

	if err := writeHeader(file); err != nil {
		s.errCh <- err
		return
	}

	res := models.Result{
		Command: "new",
		Data:    []byte(filename),
	}

	s.resCh <- res
}

func isFileExists(filename string) bool {
	fullName := filename + "." + EXTENSION
	if _, err := os.Stat(filepath.Join(SECRETS_DIR, fullName)); err == nil {
		fmt.Println("~ File with this name already exists")
		return true
	} else if !os.IsNotExist(err) {
		// В случае ошибки, отличной от "файл не существует"
		fmt.Println("~ Something went wrong with error:", err)
		return false
	}

	return false
}

func writeHeader(file *os.File) error {
	h := header()

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, h); err != nil {
		return err
	}

	if _, err := file.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func getOwnerID() ([]byte, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 {
			hashed := sha256.Sum256(iface.HardwareAddr)
			return hashed[:8], nil // Вернем первые 8 байт
		}
	}
	return nil, errors.New("no valid MAC address found")
}

func nonce() ([12]byte, error) {
	// Генерация nonce для заголовка
	var nonce [12]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return [12]byte{}, err
	}

	return nonce, nil
}

func header() models.FileHeader {
	ownerID, _ := getOwnerID()
	currTime := uint64(time.Now().Unix())
	nonce, _ := nonce()

	return models.FileHeader{
		Version:        1,
		CompleteFlag:   1, // Завершено
		OwnerID:        [8]byte(ownerID),
		SecretCount:    0,          // Количество секретов
		CreatedAt:      currTime,   // Текущая временная метка
		ModifiedAt:     currTime,   // Текущая временная метка
		DataSize:       0,          // Размер данных
		EncryptionAlgo: 0x01,       // AES-256 GCM
		Reserved:       [13]byte{}, // Заполняем резерв
		Nonce:          nonce,      // Генерируем nonce
		Checksum:       [32]byte{}, // Контрольная сумма (изначально пусто)
		Reserved2:      [32]byte{}, // Дополнительное резервное место
	}
}

func (s *secretService) Add(ctx context.Context, secret models.AddSecretCmdParams) {
	go func() {
		// zip
		// enc
		// open file
		// write header
		// write data
		// close file

		// zip.Compress(ctx, secret.Payload, s.resCh<-, s.errCh<-)
		// enc.Encrypt(ctx, s.resCh, s.errCh<-)

		jsonData, err := json.Marshal(secret)
		if err != nil {
			s.errCh <- err
			return
		}

		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)

		_, err = gzipWriter.Write(jsonData)
		if err != nil {
			s.errCh <- err
			return
		}

		if err := gzipWriter.Close(); err != nil {
			s.errCh <- err
			return
		}

		//s.resCh <- buf.Bytes()
	}()

}
