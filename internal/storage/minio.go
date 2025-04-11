package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
type MinioService interface{}

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
