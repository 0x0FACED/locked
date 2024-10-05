package zip

import (
	"context"

	"github.com/0x0FACED/locked/internal/core/models/types"
)

type Decompressor interface {
	Decompress(ctx context.Context, data []types.SecretPayload, resultCh chan<- []byte, errCh chan<- error) // Разжатие данных
}
