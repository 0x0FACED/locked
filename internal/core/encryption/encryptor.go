package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type Encryptor interface {
	Encrypt(ctx context.Context, data []byte) ([]byte, error)
}

type aesEncryptor struct {
	key   []byte   // хэш мастер-пароля
	nonce [12]byte // nonce, уникальный для каждого ФАЙЛА (в заголовке лежит)
}

func NewAesEncryptor(key []byte, nonce [12]byte) *aesEncryptor {
	return &aesEncryptor{key: key, nonce: nonce}
}

func (e *aesEncryptor) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, e.nonce[:], data, nil)
	return ciphertext, nil
}
