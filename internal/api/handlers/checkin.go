package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"io"
	"je-suis-ici-activitypub/internal/services"
	"net/http"
)

// CheckinHandler handle checkin requests
type CheckinHandler struct {
	checkinService services.CheckinService
	mediaService   services.MediaService
	serverHost     string
}

// NewCheckinHandler
func NewCheckinHandler(checkinService services.CheckinService, mediaService services.MediaService, serverHost string) *CheckinHandler {
	return &CheckinHandler{
		checkinService: checkinService,
		mediaService:   mediaService,
		serverHost:     serverHost,
	}
}

// RegisterCheckinRoutes register checkin handler routes
func (ch *CheckinHandler) RegisterCheckinRoutes(r chi.Router) {
	r.Post("/media", ch.UploadMedia)
	r.Post("/checkin", ch.CreateCheckin)
}

// CreateCheckin
func (ch *CheckinHandler) CreateCheckin(w http.ResponseWriter, r *http.Request) {
	// get user id
	userIDFromRequest, err := GetUserIDFromRequest(r)
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
	_, err := GetUserIDFromRequest(r)
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

// GetUserIDFromRequest
// valid user
// get user id from the request with JWT token
func GetUserIDFromRequest(r *http.Request) (string, error) {
	// get JWT token claims (claims is a map)
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return "", fmt.Errorf("unauthorized: %w", err)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}

	return userID, nil
}
