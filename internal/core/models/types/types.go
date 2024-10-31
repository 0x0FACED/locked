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
	Unknown    SecretType = iota // неизвестный тип
	TextFile                     // текст, например txt, csv и т.д.
	Text                         // простой текст
	Document                     // документы, например .md, .docx, .pdf и т.д.
	Archive                      // архивы, например .zip, .tar, .rar и т.д.
	Video                        // видео, например .mp4, .avi и т.д.
	Image                        // изображения
	Audio                        // аудио
	Executable                   // исполняемые файлы, например .exe, .bin и т.д.
)
