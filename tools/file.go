package tools

import (
	"errors"
	"github.com/h2non/filetype"
	"mime/multipart"
)

var (
	UnknownFileType = errors.New("Unknown file type")
)

func GetExt(file multipart.File) (string, error) {
	b := make([]byte, 262)
	if _, err := file.Read(b); err != nil {
		return "", err
	}
	kind, _ := filetype.Match(b)
	if kind == filetype.Unknown {
		return "", UnknownFileType
	}

	return kind.Extension, nil
}
