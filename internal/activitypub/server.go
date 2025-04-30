package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"je-suis-ici-activitypub/internal/db/models"
	"time"
)

type ActivityPubRepository interface {
	SaveActivity(ctx context.Context, activityID, actor, activityType, objectID, objectType, target string, rawContent []byte) error
	GetUserInboxActivities(ctx context.Context, userID uuid.UUID) ([]Activity, error)
	GetUnprocessedActivities(ctx context.Context, limit int) ([]Activity, error)
	MarkActivityAsProcessed(ctx context.Context, activityID string) error
}

type ActivityPubRepositoryImplement struct {
	pool *pgxpool.Pool
}

func NewActivityPubRepository(pool *pgxpool.Pool) ActivityPubRepository {
	return &ActivityPubRepositoryImplement{pool: pool}
}

// SaveActivity
func (apr *ActivityPubRepositoryImplement) SaveActivity(ctx context.Context, activityID, actor, activityType, objectID, objectType, target string, rawContent []byte) error {
	// default processed is false
	query := `
		INSERT INTO activities(
			activity_id, actor, type, object_id, object_type, target, raw_content, processed
		) VALUES ($1, $2, $3, $4, $5, $6, $7, false)
	`

	_, err := apr.pool.Exec(ctx, query, activityID, actor, activityType, objectID, objectType, target, rawContent)

	if err != nil {
		return fmt.Errorf("fail to save activity: %w", err)
	}
	return nil
}

// GetUserInboxActivities retrieves activities from a user's inbox
func (apr *ActivityPubRepositoryImplement) GetUserInboxActivities(ctx context.Context, userID uuid.UUID) ([]Activity, error) {
	// Query to get activities where this user is the target
	query := `
        SELECT raw_content
        FROM activities
        WHERE target = (SELECT actor_id FROM users WHERE id = $1)
        ORDER BY created_at DESC
    `

	rows, err := apr.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query inbox activities: %w", err)
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var rawContent []byte
		err := rows.Scan(&rawContent)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}

		var activity Activity
		err = json.Unmarshal(rawContent, &activity)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal activity: %w", err)
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// GetUnprocessedActivities
func (apr *ActivityPubRepositoryImplement) GetUnprocessedActivities(ctx context.Context, limit int) ([]Activity, error) {
	query := `
		SELECT activity_id, actor, type, object_id, object_type, target, raw_content
		FROM activities
		WHERE processed = false
		ORDER BY created_at ASC 
		LIMIT $1
	`

	rows, err := apr.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("fail to get unprocessed activities: %w", err)
	}
	defer rows.Close()

	var activities []Activity

	for rows.Next() {
		var activity Activity
		var activityID, actor, activityType, objectID, objectType, target string
		var rawContent []byte

		err := rows.Scan(&activityID, &actor, &activityType, &objectID, &objectType, &target, &rawContent)
		if err != nil {
			return nil, fmt.Errorf("fail to scan activity: %w", err)
		}

		// parse raw JSON format content
		err = json.Unmarshal(rawContent, &activity)
		if err != nil {
			continue
		}

		activities = append(activities, activity)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error on iterating activity rows: %w", err)
	}

	return activities, nil
}

// MarkActivityAsProcessed
func (apr *ActivityPubRepositoryImplement) MarkActivityAsProcessed(ctx context.Context, activityID string) error {
	query := `
		UPDATE activities
		SET processed = true
		WHERE activity_id = $1
	`

	_, err := apr.pool.Exec(ctx, query, activityID)
	if err != nil {
		return fmt.Errorf("fail to mark activity as processed: %w", err)
	}

	return nil
}

// FollowerRepository manage actor's followers
type FollowerRepository interface {
	AddFollower(ctx context.Context, userID uuid.UUID, followerActorID, followerInbox string) error
	RemoveFollower(ctx context.Context, userID uuid.UUID, followerActorID string) error
	GetFollowers(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type FollowerRepositoryImplement struct {
	pool *pgxpool.Pool
}

func NewFollowerRepository(pool *pgxpool.Pool) FollowerRepository {
	return &FollowerRepositoryImplement{pool: pool}
}

// AddFollower
func (fr *FollowerRepositoryImplement) AddFollower(ctx context.Context, userID uuid.UUID, followerActorID, followerInbox string) error {
	query := `
		INSERT INTO followers(user_id, follower_actor_id, follower_inbox)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, follower_actor_id) DO NOTHING
	`

	_, err := fr.pool.Exec(ctx, query, userID, followerActorID, followerInbox)
	if err != nil {
		return fmt.Errorf("fail to add follower: %w", err)
	}

	return nil
}

// RemoveFollower
func (fr *FollowerRepositoryImplement) RemoveFollower(ctx context.Context, userID uuid.UUID, followerActorID string) error {
	query := `
DELETE FROM followers
WHERE user_id = $1 AND follower_actor_id = $2
`

	_, err := fr.pool.Exec(ctx, query, userID, followerActorID)
	if err != nil {
		return fmt.Errorf("fail to remove follower: %w", err)
	}

	return nil
}

// GetFollowers
// TODO: return []Follower
// TODO: 分批取資料
func (fr *FollowerRepositoryImplement) GetFollowers(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `
		SELECT follower_actor_id, follower_inbox
		FROM followers
		WHERE user_id = $1
	`

	rows, err := fr.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("fail to get followers: %w", err)
	}
	defer rows.Close()

	var followers []string

	for rows.Next() {
		var actorID, inbox string
		err := rows.Scan(&actorID, &inbox)
		if err != nil {
			return nil, fmt.Errorf("fail to scan follower: %w", err)
		}

		followers = append(followers, inbox)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error on iterating follower rows: %w", err)
	}

	return followers, nil
}

// ActivityPubServerService
type ActivityPubServerService struct {
	activityPubRepo ActivityPubRepository
	followerRepo    FollowerRepository
	userRepo        models.UserRepository
	checkinRepo     models.CheckinRepository
	actorService    ActorService
	clientService   ActivityPubClientService
	serverHost      string
}

func NewActivityPubServerService(
	activityPubRepo ActivityPubRepository,
	followerRepo FollowerRepository,
	userRepo models.UserRepository,
	checkinRepo models.CheckinRepository,
	actorService ActorService,
	clientService ActivityPubClientService,
	serverHost string,
) *ActivityPubServerService {
	return &ActivityPubServerService{
		activityPubRepo: activityPubRepo,
		followerRepo:    followerRepo,
		userRepo:        userRepo,
		checkinRepo:     checkinRepo,
		actorService:    actorService,
		clientService:   clientService,
		serverHost:      serverHost,
	}
}

// HandleInbox handle user inbox request
func (aps *ActivityPubServerService) HandleInbox(ctx context.Context, userID uuid.UUID, body []byte) error {
	// parse activity
	var activity Activity
	err := json.Unmarshal(body, &activity)
	if err != nil {
		return fmt.Errorf("fail to parse activity: %w", err)
	}

	// get activity information
	activityID := activity.ID
	actor := activity.Actor
	activityType := activity.Type

	// parse object information
	var objectID string
	var objectType string
	var target string

	//object := activity.Object
	switch object := activity.Object.(type) {
	case string:
		objectID = object
	case map[string]interface{}:
		id, ok := object["id"].(string)
		if ok {
			objectID = id
		}
		typ, ok := object["type"].(string)
		if ok {
			objectType = typ
		}
	}

	if activity.Target != "" {
		target = activity.Target
	}

	// save activity
	err = aps.activityPubRepo.SaveActivity(ctx, activityID, actor, activityType, objectID, objectType, target, body)
	if err != nil {
		return fmt.Errorf("fail to save activity: %w", err)
	}

	// handle activity by type
	switch activityType {
	case ActivityTypeFollow:
		return aps.handleFollowActivity(ctx, userID, actor)

	case ActivityTypeUndo:
		if objectType == ActivityTypeFollow {
			return aps.handleUndoFollowActivity(ctx, userID, actor)
		}
	}

	return nil
}

func (aps *ActivityPubServerService) handleFollowActivity(ctx context.Context, userID uuid.UUID, followerActorID string) error {
	// get follower information
	follower, err := aps.clientService.FetchActorPublicInformation(ctx, followerActorID)
	if err != nil {
		return fmt.Errorf("fail to get follower actor: %w", err)
	}

	// add as follower
	err = aps.followerRepo.AddFollower(ctx, userID, follower.ID, follower.Inbox)
	if err != nil {
		return fmt.Errorf("fail to add follower: %w", err)
	}

	// get user information
	user, err := aps.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("fail to get user: %w", err)
	}

	// create Accept activity
	accept := &Activity{
		Context: DefaultContext(),
		ID:      fmt.Sprintf("https//%s/activities/%s", aps.serverHost, uuid.New().String()),
		Type:    ActivityTypeAccept,
		Actor:   user.ActorID,
		Object: map[string]interface{}{
			"id":     followerActorID,
			"type":   ActivityTypeFollow,
			"actor":  followerActorID,
			"object": user.ActorID,
		},
		To:        []string{followerActorID},
		Published: time.Now(),
	}

	// send activity
	return aps.clientService.SendActivityToTargetInbox(ctx, accept, user, follower.Inbox)
}

func (aps *ActivityPubServerService) handleUndoFollowActivity(ctx context.Context, userID uuid.UUID, followerActorID string) error {
	return aps.followerRepo.RemoveFollower(ctx, userID, followerActorID)
}

// SendActivityToInbox sends an activity to a user's inbox
func (aps *ActivityPubServerService) SendActivityToInbox(ctx context.Context, activity *Activity, sender *models.User, targetInbox string) error {
	// Use the client service to send the activity
	return aps.clientService.SendActivityToTargetInbox(ctx, activity, sender, targetInbox)
}

// GetUserInboxActivities
func (aps *ActivityPubServerService) GetUserInboxActivities(ctx context.Context, userID uuid.UUID) ([]Activity, error) {
	return aps.activityPubRepo.GetUserInboxActivities(ctx, userID)
}
