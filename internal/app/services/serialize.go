package services

import (
	"context"

	"github.com/0x0FACED/locked/internal/core/models"
	"github.com/0x0FACED/locked/internal/core/models/types"
)

func (s *secretService) Serialize(ctx context.Context, secret models.AddSecretCmdParams) ([]byte, error) {
	var rec models.SecretRecord
	off := s.db.Offset()
	rec.Offset = types.SecretOffset(off)
	
}
