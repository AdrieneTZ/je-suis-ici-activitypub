package user

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"je-suis-ici-activitypub/internal/db/models"
	"je-suis-ici-activitypub/internal/services/activitypub"
)

// UserService
type UserService interface {
	Register(ctx context.Context, serverHost, username, email, password string) (*models.User, error)
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
