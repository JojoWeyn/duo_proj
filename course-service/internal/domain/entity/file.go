package entity

import (
	"mime/multipart"
	"strings"
)

type File struct {
	Data multipart.File
	Name string
	Size int64
	Type string
}

func NewFile(file multipart.File, fileName string, fileSize int64, fileType string) *File {
	safeFileName := strings.ReplaceAll(fileName, " ", "_")

	return &File{
		Data: file,
		Name: safeFileName,
		Size: fileSize,
		Type: fileType,
	}
}
