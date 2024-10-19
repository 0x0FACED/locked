package models

import "github.com/0x0FACED/locked/internal/core/models/types"

// Эта структура нужна для отображения секрета в консоли/web/desktop UI

type SecretUI struct {
	ID          types.SecretID
	Name        types.SecretName
	Description types.SecretDescription
	Type        types.SecretType
	CreatedAt   types.SecretCreatedAt
	Size        types.SecretSize
	// Payload как таковая не будет отображаться, если секрет - файл или длинный текст.
	Payload types.SecretPayload
}

// Это параметры секрета (то, что вводится в консоль)
// А там вводиться будут только название, описание и путь до файла (или сам секрет в виде текста)
type AddSecretCmdParams struct {
	Name        *types.SecretName
	Description *types.SecretDescription
	// Строка, ибо эта структура будет использоваться, когда пользователь вводит данные в консоль чисто.
	// Поэтому у нас полезная нагрузка - это строка.
	// Если это будет просто текст - значит это и есть секрет.
	// Однако если это будет ПУТЬ к файлу или его название...
	//
	// Думаю, что определять, просто ли текст здесь или путь до файла/название файла,
	// Будем через постфикс в конце слова
	// Условно, если в конце есть .txt или .jpg, значит это ФАЙЛ.
	// Если же будет такая запись: /path/to/file, то прога это воспримет как ТЕКСТ (в конце нет расширения файла).
	Payload string
}
