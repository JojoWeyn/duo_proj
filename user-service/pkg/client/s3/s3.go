package s3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
)

type StorageS3 struct {
	Enpoint string
	Bucket  string
	Client  *minio.Client
}

func NewS3Client(endpoint, accessKeyID, secretKey, bucket string) (*StorageS3, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return &StorageS3{
		Enpoint: endpoint,
		Bucket:  bucket,
		Client:  client,
	}, nil
}

func (s *StorageS3) UploadAvatar(ctx context.Context, avatarFile multipart.File, fileName string, fileSize int64) (string, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, avatarFile)
	if err != nil {
		return "", fmt.Errorf("failed to hash avatar file: %w", err)
	}

	fileHash := fmt.Sprintf("%x", hasher.Sum(nil))
	uniqueFileName := fmt.Sprintf("%s_%s", fileHash, fileName)

	_, err = avatarFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek avatar file: %w", err)
	}

	_, err = s.Client.PutObject(ctx, s.Bucket, uniqueFileName, avatarFile, fileSize, minio.PutObjectOptions{
		ContentType: "image/png",
		UserMetadata: map[string]string{
			"x-amz-content-sha256": fmt.Sprintf("%x", hasher.Sum(nil)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s/%s/%s", s.Enpoint, s.Bucket, uniqueFileName)
	return fileURL, nil
}
