package services

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0x0FACED/locked/internal/core/database"
	"github.com/0x0FACED/locked/internal/core/encryption"
	"github.com/0x0FACED/locked/internal/core/models"
	"github.com/0x0FACED/locked/internal/core/models/types"
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

	res := models.Result{
		Command: "open",
		Data:    []byte(f.Name()),
	}

	s.resCh <- res
}

func (s *secretService) CreateSecretFile(ctx context.Context, filename string) {
	if isSecretFileExists(filename) {
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

func isSecretFileExists(filename string) bool {
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

func isFileExists(filename string) (os.FileInfo, error) {
	// проверяем состояние файла
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// не существует
		return nil, errors.New("not exists") // урааа, говнокод!
	} else if err != nil {
		// другая ошибка (доступа, например)
		return nil, err
	}

	// удостовериться, что не директория
	if fileInfo.IsDir() {
		return nil, errors.New("not a file")
	}

	// все гуд, отдаем fileInfo
	return fileInfo, nil
}

func readFile(filename string) ([]byte, error) {
	// Читаем содержимое файла в массив байтов
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return data, nil
}

func fileType(filename string) types.SecretType {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)

	switch {
	case strings.HasPrefix(mimeType, "text/"):
		return types.TextFile
	case strings.HasPrefix(mimeType, "image/"):
		return types.Image
	case strings.HasPrefix(mimeType, "video/"):
		return types.Video
	case strings.HasPrefix(mimeType, "audio/"):
		return types.Audio
	case strings.HasPrefix(mimeType, "application/"):
		switch {
		case strings.Contains(mimeType, "pdf") ||
			strings.Contains(mimeType, "msword") ||
			strings.Contains(mimeType, "vnd.ms-excel") ||
			strings.Contains(mimeType, "vnd.openxmlformats-officedocument") ||
			strings.Contains(mimeType, "rtf") ||
			strings.Contains(mimeType, "postscript"):
			return types.Document
		case strings.Contains(mimeType, "zip") ||
			strings.Contains(mimeType, "x-tar") ||
			strings.Contains(mimeType, "x-rar-compressed") ||
			strings.Contains(mimeType, "x-7z-compressed") ||
			strings.Contains(mimeType, "x-bzip2") ||
			strings.Contains(mimeType, "x-gzip"):
			return types.Archive
		case strings.Contains(mimeType, "x-msdownload") ||
			strings.Contains(mimeType, "x-executable") ||
			strings.Contains(mimeType, "octet-stream"):
			return types.Executable
		default:
			return types.Unknown
		}
	default:
		return types.Unknown
	}
}

func (s *secretService) Add(ctx context.Context, secret models.AddSecretCmdParams) {

	type serialize struct {
		Offset    uint64
		Name      *string
		Desc      *string
		Type      types.SecretType
		CreatedAt uint64
		Size      uint64
		Payload   []byte
	}

	var ser serialize

	ser.Offset = s.db.Offset()
	ser.Name = secret.Name
	ser.Desc = secret.Description
	ser.CreatedAt = uint64(time.Now().Unix())
	if secret.IsFile {
		fInfo, err := isFileExists(*secret.Name)
		if err != nil {
			s.errCh <- err
			return
		}

		ser.Size = uint64(fInfo.Size()) // сейвим исходный размер

		data, err := readFile(*secret.Name) // получаем байтовый слайс

		if err != nil {
			s.errCh <- err
			return
		}

		ser.Payload = data
	} else {
		ser.Payload = []byte(*secret.Name)
	}
	// serialize
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

}
