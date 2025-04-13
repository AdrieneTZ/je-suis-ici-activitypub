package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"je-suis-ici-activitypub/internal/db/models"
	"je-suis-ici-activitypub/internal/storage"
)

// MediaService
type MediaService interface {
	UploadMedia(ctx context.Context, data []byte, fileType, contentType string) (*models.Media, error)
	GetMediaByID(ctx context.Context, id uuid.UUID) (*models.Media, error)
}

// MediaServiceImplement
type MediaServiceImplement struct {
	mediaRepo    models.MediaRepository
	minioService storage.MinioService
}

// NewMediaService
func NewMediaService(mediaRepo models.MediaRepository, minioService storage.MinioService) MediaService {
	return &MediaServiceImplement{
		mediaRepo:    mediaRepo,
		minioService: minioService,
	}
}

// UploadMedia
// upload media file then store file name and related information to media table
// return media data including media file URL
func (ms *MediaServiceImplement) UploadMedia(ctx context.Context, fileData []byte, fileType, contentType string) (*models.Media, error) {
	// upload media file to minio
	filePath, err := ms.minioService.UploadFile(ctx, fileData, fileType, contentType)
	if err != nil {
		return nil, fmt.Errorf("fail to upload media file: %w", err)
	}

	// build media model
	media := &models.Media{
		// initially CheckinID is nil, it will be filled after checkin data is created
		FilePath: filePath,
		FileType: fileType,
		FileSize: len(fileData),
	}

	// store media
	err = ms.mediaRepo.CreateMedia(ctx, media)
	if err != nil {
		return nil, err
	}

	// generate media file URL
	fileURL, err := ms.minioService.GetFileURL(ctx, media.FilePath)
	// only if GetFileURL success, update URL field in media model
	if err == nil {
		media.URL = fileURL
	}

	return media, nil
}

// GetMediaByID
func (ms *MediaServiceImplement) GetMediaByID(ctx context.Context, id uuid.UUID) (*models.Media, error) {
	media, err := ms.mediaRepo.GetMediaByID(ctx, id)
	if err != nil {
		return nil, err
	}

	fileURL, err := ms.minioService.GetFileURL(ctx, media.FilePath)
	if err == nil {
		media.URL = fileURL
	}

	return media, nil
}
