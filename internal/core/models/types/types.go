package types

type SecretID uint8

type SecretName [64]byte

type SecretDescription [128]byte

type SecretOffset uint64

type SecretType uint8

type SecretCreatedAt uint64

type SecretSize uint64

type SecretPayload []byte

const (
	Unknown SecretType = iota
	Text
	Txt
	Archive
	Video
	Image
)
