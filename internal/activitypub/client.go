package activitypub

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"je-suis-ici-activitypub/internal/db/models"
	"net/http"
	"time"
)

// ActivityPubClientService interact with other activitypub servers
type ActivityPubClientService interface {
	FetchActorPublicInformation(ctx context.Context, actorURL string) (*Person, error)
	SendActivityToTargetInbox(ctx context.Context, activity *Activity, user *models.User, targetInbox string) error
	GetActorInbox(ctx context.Context, actorURL string) (string, error)
	GetActorFollowers(ctx context.Context, followersURL string) ([]string, error)
}

// HTTPClient send http request and return http response
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ActivityPubClientServiceImplement struct {
	httpClient HTTPClient
}

func NewActivityPubClientService(httpClient HTTPClient) ActivityPubClientService {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &ActivityPubClientServiceImplement{
		httpClient: httpClient,
	}
}

// FetchActorPublicInformation
func (ac *ActivityPubClientServiceImplement) FetchActorPublicInformation(ctx context.Context, actorURL string) (*Person, error) {
	// create http request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, actorURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to create http request: %w", err)
	}

	// set http request header
	req.Header.Set("Accept", "application/activity+json")
	req.Header.Set("User-Agent", "je-suis-ici-activitypub")

	// send http request
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fail to send http request: %w", err)
	}
	defer resp.Body.Close()

	// check http response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receive error status: %d", resp.StatusCode)
	}

	// decode http response body
	var person Person
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&person)
	if err != nil {
		return nil, fmt.Errorf("fail to decode actor public information: %w", err)
	}

	return &person, nil
}

// SendActivityToTargetInbox
func (ac *ActivityPubClientServiceImplement) SendActivityToTargetInbox(ctx context.Context, activity *Activity, user *models.User, targetInbox string) error {
	// parse activity to json
	activityJSON, err := json.Marshal(activity)
	if err != nil {
		return fmt.Errorf("fail to parse activity to json: %w", err)
	}

	// create http request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetInbox, bytes.NewBuffer(activityJSON))
	if err != nil {
		return fmt.Errorf("fail to create http request: %w", err)
	}

	// set header
	req.Header.Set("Content-Type", "application/activity+json")
	req.Header.Set("Accept", "application/activity+json")
	req.Header.Set("User-Agent", "je-suis-ici-activitypub")

	// if the user has a private key, sign the HTTP request for authentication
	// this is crucial for ActivityPub's security model
	if user.PrivateKey != "" {
		err := ac.signRequest(req, user)
		if err != nil {
			return fmt.Errorf("fail to sign request: %w", err)
		}
	}

	// send activity request
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fail to send activity request: %w", err)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rceiver error status: %d", resp.StatusCode)
	}

	return nil
}

// GetActorInbox
func (ac *ActivityPubClientServiceImplement) GetActorInbox(ctx context.Context, actorURL string) (string, error) {
	// get actor information
	actor, err := ac.FetchActorPublicInformation(ctx, actorURL)
	if err != nil {
		return "", fmt.Errorf("fail to get actor information: %w", err)
	}

	// check actor inbox
	if actor.Inbox == "" {
		return "", fmt.Errorf("actor doesn't have an inbox")
	}

	return actor.Inbox, nil
}

// GetActorFollowers
func (ac *ActivityPubClientServiceImplement) GetActorFollowers(ctx context.Context, followersURL string) ([]string, error) {
	// create http request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, followersURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to create http request: %w", err)
	}

	// set header
	req.Header.Set("Accept", "application/activity+json")
	req.Header.Set("User-Agent", "je-suis-ici-activitypub")

	// send request
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fail to get actor followers: %w", err)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receive error status: %d", resp.StatusCode)
	}

	// parse response followers
	var rawJSON map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&rawJSON)
	if err != nil {
		return nil, fmt.Errorf("fail to decode response followers: %w", err)
	}

	// get followers list
	var followers []string

	// check if item or orderedItems exists
	if items, ok := rawJSON["item"].([]interface{}); ok {
		for _, item := range items {
			follower, ok := item.(string)
			if ok {
				followers = append(followers, follower)
			}
		}
	} else if orderedItems, ok := rawJSON["orderedItems"].([]interface{}); ok {
		for _, orderedItem := range orderedItems {
			follower, ok := orderedItem.(string)
			if ok {
				followers = append(followers, follower)
			}
		}
	} else {
		return nil, fmt.Errorf("followers collection doesn't have item or orderedItems")
	}

	return followers, nil
}

// signRequest sign an HTTP request using RSA cryptography
func (ac *ActivityPubClientServiceImplement) signRequest(req *http.Request, user *models.User) error {
	// decodes the PEM-encoded private key
	block, _ := pem.Decode([]byte(user.PrivateKey))
	if block == nil {
		return fmt.Errorf("fail to decode private key")
	}

	// parse the decoded key bytes as an RSA private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("fail to parse private key: %w", err)
	}

	// extract values needed for the signature
	method := req.Method
	path := req.URL.Path
	host := req.URL.Host
	// create a formatted UTC timestamp and set it as the Date to request header
	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set("Date", date)

	// create the string to be signed
	// follow HTTP Signature specification
	signString := fmt.Sprintf("(request-target): %s %s\nhost: %s\ndate: %s",
		method, path, host, date)

	// compute the SHA-256 hash of the string to be signed
	h := sha256.New()
	h.Write([]byte(signString))
	digest := h.Sum(nil)

	// sign the digest using the RSA private key with PKCS#1 v1.5 padding
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest)
	if err != nil {
		return fmt.Errorf("fail to sign: %w", err)
	}

	// encode the binary signature as a Base64 string
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	// create a key identifier using the user's ActivityPub Actor ID
	keyId := fmt.Sprintf("%s#main-key", user.ActorID)

	// format the HTTP Signature header with key ID, algorithm, signed headers, and the signature
	signatureHeader := fmt.Sprintf(`keyId="%s",algorithm="rsa-sha256",headers="(request-target) host date",signature="%s"`,
		keyId, encodedSignature)

	// add the signature header to the HTTP request
	req.Header.Set("Signature", signatureHeader)

	return nil
}
