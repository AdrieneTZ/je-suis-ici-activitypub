package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"je-suis-ici-activitypub/internal/activitypub"
	"je-suis-ici-activitypub/internal/db/models"
	"time"
)

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
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("fail to hash password: %w", err)
	}

	// generate actorID
	actorID := us.actorService.GenerateActorID(serverHost, username)

	// generate private and public key pair
	privateKey, publicKey, err := us.actorService.GenerateKeyPair()
	if err != nil {
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
		return nil, err
	}

	return user, nil
}

// Authenticate verify user account
func (us *UserServiceImplement) Authenticate(ctx context.Context, usernameOrEmail, password string) (*models.User, error) {
	var user *models.User
	var err error

	// get user data by username or by email
	user, err = us.userRepo.GetByUsername(ctx, usernameOrEmail)
	if err != nil {
		user, err = us.userRepo.GetByEmail(ctx, usernameOrEmail)
		if err != nil {
			return nil, fmt.Errorf("invalid credentials")
		}
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// GetUserByID
func (us *UserServiceImplement) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return us.userRepo.GetByID(ctx, id)
}

// GetUserByUsername
func (us *UserServiceImplement) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return us.userRepo.GetByUsername(ctx, username)
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
