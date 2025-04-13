package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"je-suis-ici-activitypub/internal/db/models"
	"je-suis-ici-activitypub/internal/storage"
)

// CheckinService
type CheckinService interface {
	CreateCheckin(ctx context.Context, userID uuid.UUID, content, locationName string, latitude, longitude float64, mediaIDs []uuid.UUID, serverHost string) (*models.Checkin, error)
	GetCheckinByID(ctx context.Context, id uuid.UUID) (*models.Checkin, error)
}

// CheckinServiceImplement
type CheckinServiceImplement struct {
	checkinRepo  models.CheckinRepository
	mediaRepo    models.MediaRepository
	minioService storage.MinioService
}

// NewCheckinService
func NewCheckinService(checkinRepo models.CheckinRepository, mediaRepo models.MediaRepository, minioService storage.MinioService) CheckinService {
	return &CheckinServiceImplement{
		checkinRepo:  checkinRepo,
		mediaRepo:    mediaRepo,
		minioService: minioService,
	}
}

// CreateCheckin
func (cs *CheckinServiceImplement) CreateCheckin(ctx context.Context, userID uuid.UUID, content, locationName string, latitude, longitude float64, mediaIDs []uuid.UUID, serverHost string) (*models.Checkin, error) {
	// generate ActivityPub activities ID
	activityID := fmt.Sprintf("https://%s/activities/%s", serverHost, uuid.New().String())

	// build checkin model
	checkin := &models.Checkin{
		UserID:       userID,
		Content:      content,
		LocationName: locationName,
		Latitude:     latitude,
		Longitude:    longitude,
		ActivityID:   activityID,
	}

	// store checkin
	err := cs.checkinRepo.CreateCheckin(ctx, checkin)
	if err != nil {
		return nil, err
	}

	// create relation between checkin and media
	if len(mediaIDs) > 0 {
		for _, mediaID := range mediaIDs {
			media, err := cs.mediaRepo.GetMediaByID(ctx, mediaID)
			if err != nil {
				continue
			}

			// check if media data has related to checkin data
			if media.CheckinID != uuid.Nil {
				continue
			}
			// create relation between media data and checkin data by set FK to media data
			media.CheckinID = checkin.ID

			// update media data
			err = cs.mediaRepo.UpdateMedia(ctx, media)
			if err != nil {
				continue
			}

			checkin.Media = append(checkin.Media, *media)
		}
	}

	// get full checkin data
	fullCheckin, err := cs.checkinRepo.GetCheckinByID(ctx, checkin.ID)
	if err != nil {
		return checkin, nil
	}

	return fullCheckin, nil
}

// GetCheckinByID
func (cs *CheckinServiceImplement) GetCheckinByID(ctx context.Context, id uuid.UUID) (*models.Checkin, error) {
	checkin, err := cs.checkinRepo.GetCheckinByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to get checkin: %w", err)
	}

	// generate media file URL
	for i := range checkin.Media {
		url, err := cs.minioService.GetFileURL(ctx, checkin.Media[i].FilePath)
		if err == nil {
			checkin.Media[i].URL = url
		}
	}

	return checkin, nil
}
