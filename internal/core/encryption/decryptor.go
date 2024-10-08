package encryption

import (
	"context"

	"github.com/0x0FACED/locked/internal/core/models/types"
)

type Decryptor interface {
	Decrypt(ctx context.Context, data []types.SecretPayload, resultCh chan<- []byte, errCh chan<- error)
}
