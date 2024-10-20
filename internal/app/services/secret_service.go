package services

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"

	"github.com/0x0FACED/locked/internal/core/database"
	"github.com/0x0FACED/locked/internal/core/encryption"
	"github.com/0x0FACED/locked/internal/core/models"
)

type SecretService interface {
	Add(ctx context.Context, secret models.AddSecretCmdParams)
}

type secretService struct {
	zip   zip.Compressor
	unzip zip.Decompressor

	enc encryption.Encryptor
	dec encryption.Decryptor

	db database.Database

	resCh chan []byte
	errCh chan error
	done  chan struct{}
}

func New(resCh chan []byte, errCh chan error, done chan struct{}) SecretService {
	return &secretService{
		resCh: resCh,
		errCh: errCh,
		done:  done,
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

		s.resCh <- buf.Bytes()
	}()

}
