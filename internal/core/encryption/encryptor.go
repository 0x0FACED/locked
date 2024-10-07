package encryption

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/0x0FACED/locked/internal/core/models/types"
	"golang.org/x/crypto/scrypt"
)

type Encryptor interface {
	Encrypt(ctx context.Context, data []types.SecretPayload, resultCh chan<- []byte, errCh chan<- error)
}

type aesEncryptor struct {
	key []byte
}

func NewEnc(masterPassword string) (*aesEncryptor, error) {
	salt := []byte("test_salt")
	key, err := scrypt.Key([]byte(masterPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return &aesEncryptor{key: key}, nil
}

func (e *aesEncryptor) Encrypt(ctx context.Context, data []types.SecretPayload, resultCh chan<- []byte, errCh chan<- error) {
	for _, payload := range data {
		iv := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			errCh <- err
			return
		}

		block, err := aes.NewCipher(e.key)
		if err != nil {
			errCh <- err
			return
		}

		paddedData := pad(payload)

		encrypted := make([]byte, len(paddedData))
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(encrypted, paddedData)

		resultCh <- append(iv, encrypted...)
	}
}

func pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}
