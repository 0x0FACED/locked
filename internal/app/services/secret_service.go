package services

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/0x0FACED/locked/internal/core/database"
	"github.com/0x0FACED/locked/internal/core/encryption"
	"github.com/0x0FACED/locked/internal/core/models"
)

type SecretService struct {
	zip   zip.Compressor
	unzip zip.Decompressor

	enc encryption.Encryptor
	dec encryption.Decryptor

	db database.Database

	resCh chan []byte
	errCh chan error
	done  chan struct{}
}

func New() *SecretService {
	return &SecretService{}
}

// TODO: Заменить secretUI на фактические данные, которые сохраняются
func (s *SecretService) Add(secret models.SecretUI) {
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

	s.resCh <- buf.Bytes()
}
