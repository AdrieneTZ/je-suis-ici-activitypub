package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"je-suis-ici-activitypub/internal/activitypub"
	"je-suis-ici-activitypub/internal/services"
	"net/http"
	"strconv"
	"time"
)

// CheckinHandler handle checkin requests
type CheckinHandler struct {
	userService     services.UserService
	checkinService  services.CheckinService
	mediaService    services.MediaService
	apServerService *activitypub.ActivityPubServerService
	authHandler     AuthHandler
	serverHost      string
}

// NewCheckinHandler
func NewCheckinHandler(userService services.UserService, checkinService services.CheckinService, mediaService services.MediaService, apServerService *activitypub.ActivityPubServerService, authHandler AuthHandler, serverHost string) *CheckinHandler {
	return &CheckinHandler{
		userService:     userService,
		checkinService:  checkinService,
		mediaService:    mediaService,
		apServerService: apServerService,
		authHandler:     authHandler,
		serverHost:      serverHost,
	}
}

// RegisterCheckinRoutes register checkin handler routes
func (ch *CheckinHandler) RegisterCheckinRoutes(r chi.Router) {
	r.Post("/media", ch.UploadMedia)
	r.Post("/checkins", ch.CreateCheckin)
	r.Get("/checkins", ch.GetUserCheckins)
	r.Get("/checkins/{id}", ch.GetCheckinByID)
}

// CreateCheckin
func (ch *CheckinHandler) CreateCheckin(w http.ResponseWriter, r *http.Request) {
	// get user id
	userIDFromRequest, err := ch.authHandler.GetUserIDByAuthTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// turn string to uuid
	userID, err := uuid.Parse(userIDFromRequest)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
	}

	// parse request
	var req struct {
		Content      string      `json:"content"`
		LocationName string      `json:"location_name"`
		Latitude     float64     `json:"latitude"`
		Longitude    float64     `json:"longitude"`
		MediaIDs     []uuid.UUID `json:"media_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// create new checkin
	checkin, err := ch.checkinService.CreateCheckin(
		r.Context(),
		userID, req.Content, req.LocationName,
		req.Latitude, req.Longitude, req.MediaIDs, r.Host,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: public this new checkin to ActivityPub global newsfeed

	// return created checkin
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(checkin)
}

// UploadMedia
func (ch *CheckinHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	// valid user
	_, err := ch.authHandler.GetUserIDByAuthTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// parse form data from request
	err = r.ParseMultipartForm(32 << 20) // max 32 MB
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	// get upload file data from request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "no file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "fail to read file", http.StatusInternalServerError)
		return
	}

	// get file content type
	contentType := fileHeader.Header.Get("Content-Type")
	// if content type is empty, set default content type
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// upload file
	media, err := ch.mediaService.UploadMedia(r.Context(), fileData, "image", contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(media)
}

// GetCheckinByID
func (ch *CheckinHandler) GetCheckinByID(w http.ResponseWriter, r *http.Request) {
	// get checkin id
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid checkin id", http.StatusBadRequest)
	}

	// get checkin data
	checkin, err := ch.checkinService.GetCheckinByID(r.Context(), id)
	if err != nil {
		http.Error(w, "checkin not found", http.StatusNotFound)
		return
	}

	// return checkin data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkin)
}

// GetUserCheckins
func (ch *CheckinHandler) GetUserCheckins(w http.ResponseWriter, r *http.Request) {
	// get user id
	userIDFromRequest, err := ch.authHandler.GetUserIDByAuthTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// turn string to uuid
	userID, err := uuid.Parse(userIDFromRequest)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// get pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}

	// get user checkins
	checkins, err := ch.checkinService.GetCheckinsByUserID(r.Context(), userID, pageSize, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return user checkins
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"checkins":  checkins,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetUserActivityPubInfo return a user's ActivityPub information
func (ch *CheckinHandler) GetUserActivityPubInfo(w http.ResponseWriter, r *http.Request) {
	// get username from request URL
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	// get user by username
	user, err := ch.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	actorID := user.ActorID

	// return user's ActivityPub information

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"actor_id":  actorID,
		"inbox":     fmt.Sprintf("%s/inbox", actorID),
		"outbox":    fmt.Sprintf("%s/outbox", actorID),
		"followers": fmt.Sprintf("%s/followers", actorID),
	})
}

// SendCheckinToUser send a checkin activity to target user's inbox
func (ch *CheckinHandler) SendCheckinToUser(w http.ResponseWriter, r *http.Request) {
	// get authenticated user ID
	senderUserIDFromRequest, err := ch.authHandler.GetUserIDByAuthTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// turn sender user id to uuid
	senderUserID, err := uuid.Parse(senderUserIDFromRequest)
	if err != nil {
		http.Error(w, "invalid sender user id", http.StatusBadRequest)
		return
	}

	// get sender username from request URL
	senderUsername := chi.URLParam(r, "sender_username")
	if senderUsername == "" {
		http.Error(w, "sender username is required", http.StatusBadRequest)
		return
	}

	// verify the authenticated user matches the sender
	sender, err := ch.userService.GetUserByID(r.Context(), senderUserID)
	if err != nil {
		http.Error(w, "sender not found", http.StatusNotFound)
		return
	}

	// check sender user is current user
	if sender.Username != senderUsername {
		http.Error(w, "sender can only send activities as himself", http.StatusForbidden)
		return
	}

	// parse request
	var req struct {
		RecipientUsername string  `json:"recipient_username"`
		Content           string  `json:"content"`
		LocationName      string  `json:"location_name"`
		Latitude          float64 `json:"latitude"`
		Longitude         float64 `json:"longitude"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// get recipient by username
	recipient, err := ch.userService.GetUserByUsername(r.Context(), req.RecipientUsername)
	if err != nil {
		http.Error(w, "Recipient not found", http.StatusNotFound)
		return
	}

	// create a checkin object
	checkinID := uuid.New()
	checkinURL := fmt.Sprintf("https://%s/checkins/%s", ch.serverHost, checkinID)

	// create activityPub Note object for the checkin
	note := &activitypub.Object{
		Context:      activitypub.DefaultContext(),
		ID:           checkinURL,
		Type:         "Note",
		Content:      req.Content,
		Published:    time.Now().UTC(),
		AttributedTo: sender.ActorID,
		Location: &activitypub.Place{
			Type:      "Place",
			Name:      req.LocationName,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
	}

	// create activity
	activity := &activitypub.Activity{
		Context:   activitypub.DefaultContext(),
		ID:        fmt.Sprintf("%s/activity", checkinURL),
		Type:      "Create",
		Actor:     sender.ActorID,
		Object:    note,
		Published: time.Now().UTC(),
		To:        []string{recipient.ActorID},
	}

	// get recipient's inbox URL
	recipientInbox := fmt.Sprintf("%s/inbox", recipient.ActorID)

	// send activity to recipient's inbox
	err = ch.apServerService.SendActivityToInbox(r.Context(), activity, sender, recipientInbox)
	if err != nil {
		http.Error(w, fmt.Sprintf("fail to send activity: %v", err), http.StatusInternalServerError)
		return
	}

	// return success response
	resp := struct {
		Success  bool                  `json:"success"`
		Message  string                `json:"message"`
		Activity *activitypub.Activity `json:"activity"`
	}{
		Success:  true,
		Message:  fmt.Sprintf("check-in sent to %s's inbox", req.RecipientUsername),
		Activity: activity,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetUserInbox retrieves activities from a user's inbox
func (ch *CheckinHandler) GetUserInbox(w http.ResponseWriter, r *http.Request) {
	// Get username from URL
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get authenticated user ID
	userIDFromRequest, err := ch.authHandler.GetUserIDByAuthTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get user by username
	user, err := ch.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify the authenticated user matches the requested inbox owner
	if user.ID.String() != userIDFromRequest {
		http.Error(w, "You can only view your own inbox", http.StatusForbidden)
		return
	}

	// Get activities from inbox
	activities, err := ch.apServerService.GetUserInboxActivities(r.Context(), user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get inbox activities: %v", err), http.StatusInternalServerError)
		return
	}

	// Return activities
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}
