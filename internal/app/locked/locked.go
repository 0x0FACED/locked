package locked

import (
	"github.com/0x0FACED/locked/internal/core/database"
	"github.com/0x0FACED/locked/internal/core/encryption"
	"github.com/0x0FACED/locked/internal/core/zip"
)

type app struct {
	zip   zip.Compressor
	unzip zip.Decompressor

	enc encryption.Encryptor
	dec encryption.Decryptor

	db database.Database
}
