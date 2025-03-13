package s3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/JojoWeyn/duo-proj/user-service/pkg/client/s3"
	"github.com/minio/minio-go/v7"
	"io"
	"mime/multipart"
)

type UserS3Repo struct {
	StorageS3 *s3.StorageS3
}

func NewUserS3Repository(storageS3 *s3.StorageS3) *UserS3Repo {
	return &UserS3Repo{
		StorageS3: storageS3,
	}
}

func (s *UserS3Repo) UploadAvatar(ctx context.Context, avatarFile multipart.File, fileName string, fileSize int64) (string, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, avatarFile)
	if err != nil {
		return "", fmt.Errorf("failed to hash avatar file: %w", err)
	}

	_, err = avatarFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek avatar file: %w", err)
	}

	_, err = s.StorageS3.Client.PutObject(ctx, s.StorageS3.Bucket, fileName, avatarFile, fileSize, minio.PutObjectOptions{
		ContentType: "image/png",
		UserMetadata: map[string]string{
			"x-amz-content-sha256": fmt.Sprintf("%x", hasher.Sum(nil)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s/%s/%s", s.StorageS3.Enpoint, s.StorageS3.Bucket, fileName)
	return fileURL, nil
}
