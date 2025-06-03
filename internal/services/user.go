package services

import (
	"context"
	"fmt"
	"je-suis-ici-activitypub/internal/activitypub"
	"je-suis-ici-activitypub/internal/db/models"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
)

var tracer = otel.Tracer("services/user")

// UserService
type UserService interface {
	Register(ctx context.Context, serverHost, username, email, password string) (*models.User, error)
	Authenticate(ctx context.Context, usernameOrEmail, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// UserServiceImplement
type UserServiceImplement struct {
	userRepo     models.UserRepository
	actorService activitypub.ActorService
}

// NewUserService
func NewUserService(userRepo models.UserRepository, actorService activitypub.ActorService) UserService {
	return &UserServiceImplement{
		userRepo:     userRepo,
		actorService: actorService,
	}
}

// Register user register an account
func (us *UserServiceImplement) Register(ctx context.Context, serverHost, username, email, password string) (*models.User, error) {
	// add child tracer
	ctx, span := tracer.Start(ctx, "Register")
	defer span.End()

	// record input attributes
	span.SetAttributes(
		attribute.String("username", username),
		attribute.String("email", email),
	)

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("error.type", "hash_password_error"),
			attribute.String("error.message", err.Error()),
		)

		return nil, fmt.Errorf("fail to hash password: %w", err)
	}

	// generate actorID
	actorID := us.actorService.GenerateActorID(serverHost, username)

	// generate private and public key pair
	privateKey, publicKey, err := us.actorService.GenerateKeyPair()
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("error.type", "generate_key_pair_error"),
			attribute.String("error.message", err.Error()),
		)

		return nil, err
	}

	// build user model
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		DisplayName:  username,
		ActorID:      actorID,
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
	}

	// create user's account
	err = us.userRepo.CreateUser(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("error.type", "create_user_account_error"),
			attribute.String("error.message", err.Error()),
		)

		return nil, err
	}

	return user, nil
}

// Authenticate verify user account
func (us *UserServiceImplement) Authenticate(ctx context.Context, usernameOrEmail, password string) (*models.User, error) {
	// add child tracer
	ctx, span := tracer.Start(ctx, "Authenticate")
	defer span.End()

	// record input attributes
	span.SetAttributes(attribute.String("username or email", usernameOrEmail))

	var user *models.User
	var err error

	// get user data by username or by email
	user, err = us.userRepo.GetByUsername(ctx, usernameOrEmail)
	if err != nil {
		user, err = us.userRepo.GetByEmail(ctx, usernameOrEmail)
		if err != nil {
			span.RecordError(err)
			span.SetAttributes(
				attribute.String("error.type", "get_user_by_username_or_email_error"),
				attribute.String("error.message", err.Error()),
			)
			span.SetStatus(codes.Error, "authentication failed")

			return nil, fmt.Errorf("invalid credentials")
		}
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("error.type", "compare_password_error"),
			attribute.String("error.message", err.Error()),
		)

		return nil, fmt.Errorf("invalid credentials")
	}

	// record output attributes
	span.SetAttributes(
		attribute.String("user.id", user.ID.String()),
	)

	return user, nil
}

// GetUserByID
func (us *UserServiceImplement) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	// add child tracer
	ctx, span := tracer.Start(ctx, "GetUserByID")
	defer span.End()

	// record input attributes
	span.SetAttributes(attribute.String("userID", id.String()))

	user, err := us.userRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("error.type", "get_user_by_id_error"),
			attribute.String("error.message", err.Error()),
		)
		span.SetStatus(codes.Error, "fail to get user by id")

		return nil, err
	}

	// record output attributes
	if user != nil {
		span.SetAttributes(attribute.String("user.id", user.ID.String()))
	}

	return user, nil
}

// GetUserByUsername
func (us *UserServiceImplement) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	// add child tracer
	ctx, span := tracer.Start(ctx, "GetUserByUsername")
	defer span.End()

	// record input attributes
	span.SetAttributes(attribute.String("username", username))

	user, err := us.userRepo.GetByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("error.type", "get_user_by_username_error"),
			attribute.String("error.message", err.Error()),
		)
		span.SetStatus(codes.Error, "fail to get user by username")

		return nil, err
	}

	// record output attributes
	if user != nil {
		span.SetAttributes(attribute.String("user.id", user.ID.String()))
	}

	return user, nil
}

// UpdateUser
func (us *UserServiceImplement) UpdateUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return fmt.Errorf("invalid user input")
	}
	if user.ID == uuid.Nil {
		return fmt.Errorf("invalid user input")
	}

	user.UpdatedAt = time.Now()

	return us.userRepo.UpdateUser(ctx, user)
}

// DeleteUser
func (us *UserServiceImplement) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("invalid user input")
	}

	return us.userRepo.DeleteUser(ctx, id)
}
