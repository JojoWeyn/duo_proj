package s3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/course-service/pkg/client/s3"
	"github.com/minio/minio-go/v7"
	"io"
)

type FileS3Repo struct {
	StorageS3 *s3.StorageS3
}

func NewFileS3Repository(storageS3 *s3.StorageS3) *FileS3Repo {
	return &FileS3Repo{
		StorageS3: storageS3,
	}
}

func (s *FileS3Repo) UploadFile(ctx context.Context, File *entity.File) (string, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, File.Data)
	if err != nil {
		return "", fmt.Errorf("failed to hash avatar file: %w", err)
	}

	_, err = File.Data.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek avatar file: %w", err)
	}

	_, err = s.StorageS3.Client.PutObject(ctx, s.StorageS3.Bucket, File.Name, File.Data, File.Size, minio.PutObjectOptions{
		ContentType: File.Type,
		UserMetadata: map[string]string{
			"x-amz-content-sha256": fmt.Sprintf("%x", hasher.Sum(nil)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s/%s/%s", s.StorageS3.Endpoint, s.StorageS3.Bucket, File.Name)
	return fileURL, nil
}

func (s *FileS3Repo) ListFiles(ctx context.Context) ([]string, error) {
	var fileURLs []string

	objectCh := s.StorageS3.Client.ListObjects(ctx, s.StorageS3.Bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects: %w", object.Err)
		}

		fileURL := fmt.Sprintf("https://%s/%s/%s", s.StorageS3.Endpoint, s.StorageS3.Bucket, object.Key)
		fileURLs = append(fileURLs, fileURL)
	}

	return fileURLs, nil
}

func (s *FileS3Repo) DeleteFile(ctx context.Context, fileName string) error {
	return s.StorageS3.Client.RemoveObject(ctx, s.StorageS3.Bucket, fileName, minio.RemoveObjectOptions{})
}
