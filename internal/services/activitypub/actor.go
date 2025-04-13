package activitypub

import (
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
	//CreateActor(ctx context.Context, user *models.User, serverHost string) error
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

	actorId := fmt.Sprintf("https://%s/users/%s", serverHost, safeUsername)

	return actorId
}

// CreateActor create user's ActivityPub Actor
//func (as *ActorServiceImplement) CreateActor(ctx context.Context, user *models.User, serverHost string) error {
//	// generate private and public key pair
//	privateKey, publicKey, err := as.GenerateKeyPair()
//	if err != nil {
//		return fmt.Errorf("fail to generate private and public key pair: %w", err)
//	}
//
//	// encode username to make sure username is safe
//	safeUsername := url.PathEscape(user.Username)
//
//	// set actor id
//	if user.ActorID == "" {
//		user.ActorID = fmt.Sprintf("https://%s/users/%s", serverHost, safeUsername)
//	}
//
//	user.PrivateKey = privateKey
//	user.PublicKey = publicKey
//
//	return as.userRepo.Update(ctx, user)
//}
