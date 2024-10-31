package services

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/0x0FACED/locked/internal/core/models"
)

func Serialize(ctx context.Context, record models.SecretRecord) ([]byte, error) {
	buf := new(bytes.Buffer)

	// offset -> id -> name -> desc -> type -> created at -> size -> pload
	if err := binary.Write(buf, binary.LittleEndian, record.Offset); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, record.ID); err != nil {
		return nil, err
	}
	if _, err := buf.Write(record.Name[:]); err != nil {
		return nil, err
	}
	if _, err := buf.Write(record.Description[:]); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, record.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, record.CreatedAt); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, record.Size); err != nil {
		return nil, err
	}
	if _, err := buf.Write(record.Payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
