package usecase

import (
	"context"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"mime/multipart"
)

type FileS3Repository interface {
	UploadFile(ctx context.Context, File *entity.File) (string, error)
	ListFiles(ctx context.Context) ([]string, error)
	DeleteFile(ctx context.Context, fileName string) error
}

type fileS3UseCase struct {
	fileS3Repo FileS3Repository
}

func NewFileS3UseCase(fileS3Repo FileS3Repository) *fileS3UseCase {
	return &fileS3UseCase{
		fileS3Repo: fileS3Repo,
	}
}

func (f *fileS3UseCase) UploadFile(ctx context.Context, file multipart.File, fileName string, fileSize int64, fileType string) (string, error) {
	fileEntity := entity.NewFile(file, fileName, fileSize, fileType)
	return f.fileS3Repo.UploadFile(ctx, fileEntity)
}

func (f *fileS3UseCase) ListFiles(ctx context.Context) ([]string, error) {
	return f.fileS3Repo.ListFiles(ctx)
}

func (f *fileS3UseCase) DeleteFile(ctx context.Context, fileName string) error {
	return f.fileS3Repo.DeleteFile(ctx, fileName)
}
