package activitypub

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"je-suis-ici-activitypub/internal/db/models"
	"net/url"
)

type ActorService interface {
	GenerateKeyPair() (string, string, error)
	GenerateActorID(serverHost, username string) string
	CreateActor(ctx context.Context, user *models.User, serverHost string) error
	GetActor(ctx context.Context, user *models.User, serverHost string) (*Person, error)
}

type ActorServiceImplement struct {
	userRepo models.UserRepository
}

func NewActorService(userRepo models.UserRepository) ActorService {
	return &ActorServiceImplement{userRepo: userRepo}
}

// GenerateKeyPair generate private and public key pair
func (as *ActorServiceImplement) GenerateKeyPair() (string, string, error) {
	// generate 2048 bits RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("fail to generate RSA key: %w", err)
	}

	// serialize private key to PEM format
	// PEM(Privacy Enhanced Mail)
	// store PEM format to database
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// serialize public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("fail to marshal public key: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil

}

func (as *ActorServiceImplement) GenerateActorID(serverHost, username string) string {
	// encode username to make sure username is safe
	safeUsername := url.PathEscape(username)

	actorId := fmt.Sprintf("http://%s/users/%s", serverHost, safeUsername)

	return actorId
}

// CreateActor create user's ActivityPub Actor
func (as *ActorServiceImplement) CreateActor(ctx context.Context, user *models.User, serverHost string) error {
	// generate private and public key pair
	privateKey, publicKey, err := as.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("fail to generate private and public key pair: %w", err)
	}

	// encode username to make sure username is safe
	safeUsername := url.PathEscape(user.Username)

	// set actor id
	if user.ActorID == "" {
		user.ActorID = fmt.Sprintf("http://%s/users/%s", serverHost, safeUsername)
	}

	user.PrivateKey = privateKey
	user.PublicKey = publicKey

	err = as.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// GetActor
func (as *ActorServiceImplement) GetActor(ctx context.Context, user *models.User, serverHost string) (*Person, error) {
	actorID := fmt.Sprintf("http://%s/user/%s", serverHost, user.Username)

	actor := &Person{
		Context:           DefaultContext(),
		ID:                actorID,
		Type:              "Person",
		Name:              user.DisplayName,
		PreferredUsername: user.Username,
		Inbox:             fmt.Sprintf("%s/inbox", actorID),
		Outbox:            fmt.Sprintf("%s/outbox", actorID),
		Following:         fmt.Sprintf("%s/following", actorID),
		Followers:         fmt.Sprintf("%s/followers", actorID),
		Liked:             fmt.Sprintf("%s/linked", actorID),
		URL:               actorID,
		Published:         user.CreatedAt,
		Updated:           user.UpdatedAt,
	}

	if user.AvatarURL != "" {
		actor.Icon = &Image{
			Type: "Image",
			URL:  user.AvatarURL,
		}
	}

	if user.PublicKey != "" {
		actor.PublicKey = PublicKey{
			ID:           fmt.Sprintf("%s#main-key", actorID),
			Owner:        actorID,
			PublicKeyPem: user.PublicKey,
		}
	}

	return actor, nil
}
