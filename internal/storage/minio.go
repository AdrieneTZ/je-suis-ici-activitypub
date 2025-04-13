package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"
	"time"
)

// define file type
const (
	FileTypeImage = "image"
)

// MinioConfig MinIO configuration
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// MinioService
type MinioService interface {
	UploadFile(ctx context.Context, fileData []byte, fileType, contentType string) (string, error)
	GetFileURL(ctx context.Context, fileName string) (string, error)
}

// MinioServiceImplement
type MinioServiceImplement struct {
	client     *minio.Client
	bucket     string
	endpoint   string
	publicURLs bool
}

func NewMinioServiceImplement(minCfg MinioConfig) (*MinioServiceImplement, error) {
	// init MinIO client
	client, err := minio.New(minCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minCfg.AccessKey, minCfg.SecretKey, ""),
		Secure: minCfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to init MinIO client: %w", err)
	}

	// check if bucket exists
	bucketExistes, err := client.BucketExists(context.Background(), minCfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("fail to check if bucket exists: %w", err)
	}

	// if bucket doesn't exist, create one
	if !bucketExistes {
		err = client.MakeBucket(context.Background(), minCfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("fail to create new bucket: %w", err)
		}

		// set bucket as public access
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`
		policy = fmt.Sprintf(policy, minCfg.Bucket)

		err = client.SetBucketPolicy(context.Background(), minCfg.Bucket, policy)
		if err != nil {
			return nil, fmt.Errorf("fail to set bucket policy: %w", err)
		}
	}

	return &MinioServiceImplement{
		client:     client,
		bucket:     minCfg.Bucket,
		endpoint:   minCfg.Endpoint,
		publicURLs: true,
	}, nil
}

// UploadFile
func (mis *MinioServiceImplement) UploadFile(ctx context.Context, fileData []byte, fileType, contentType string) (string, error) {
	// generate unique file path (but not absolute path, it's also file name)
	filePath := fmt.Sprintf("%s/%s%s", fileType, uuid.New().String(), getExtension(contentType))

	// upload file
	_, err := mis.client.PutObject(ctx, mis.bucket, filePath, bytes.NewReader(fileData), int64(len(fileData)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("fail to upload file: %w", err)
	}

	return filePath, nil
}

// GetFileURL return public file URL
func (mis *MinioServiceImplement) GetFileURL(ctx context.Context, filePath string) (string, error) {
	if mis.publicURLs {
		// use presigned URL to allow temporarily access
		reqParams := make(url.Values)
		presignedURL, err := mis.client.PresignedGetObject(ctx, mis.bucket, filePath, time.Hour*24, reqParams)
		if err != nil {
			return "", fmt.Errorf("fail to generate presigned URL: %w", err)
		}

		return presignedURL.String(), nil
	}

	directFileURL := fmt.Sprintf("http://%s/%s/%s", mis.endpoint, mis.bucket, filePath)
	return directFileURL, nil
}

// getExtension return file type based on contentType
func getExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}
