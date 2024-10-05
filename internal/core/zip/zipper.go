package zip

import (
	"context"

	"github.com/0x0FACED/locked/internal/core/models/types"
)

type Compressor interface {
	Compress(ctx context.Context, data []types.SecretPayload, resultCh chan<- []byte, errCh chan<- error) // Сжатие данных
}
