package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"je-suis-ici-activitypub/internal/activitypub"
	"je-suis-ici-activitypub/internal/services"
	"net/http"
	"time"
)

type UserHandler struct {
	userService  services.UserService
	actorService activitypub.ActorService
	authHandler  AuthHandler
	serverHost   string
}

func NewUserHandler(userService services.UserService, actorService activitypub.ActorService, authHandler AuthHandler, serverHost string) *UserHandler {
	return &UserHandler{
		userService:  userService,
		actorService: actorService,
		authHandler:  authHandler,
		serverHost:   serverHost,
	}
}

func (uh *UserHandler) RegisterUserRouters(r chi.Router) {
	r.Put("/users/{id}", uh.UpdateUser)
	r.Delete("/users/{id}", uh.DeleteUser)
}

type UpdateUserRequest struct {
	Username    *string `json:"username,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	PublicKey   *string `json:"public-key,omitempty"`
	PrivateKey  *string `json:"private-key,omitempty"`
}

type UpdateUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name,omitempty"`
	Email       string    `json:"email"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	ActorID     string    `json:"actor_id,omitempty"`
	PublicKey   string    `json:"public_key,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// get user id from request URL params
	idFromParam := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idFromParam)
	if err != nil {
		http.Error(w, "invalid user ID format from request parameter", http.StatusBadRequest)
		return
	}

	// get authenticated user from context (set by JWT middleware)
	userIDFromToken, err := uh.authHandler.GetUserIDByAuthTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if idFromParam != userIDFromToken {
		http.Error(w, "user can only update their own profile", http.StatusForbidden)
		return
	}

	// get user data from db
	currentUser, err := uh.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updatedUser UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// add updated data
	if updatedUser.Username != nil && *updatedUser.Username != "" {
		currentUser.Username = *updatedUser.Username
		currentUser.ActorID = uh.actorService.GenerateActorID(uh.serverHost, *updatedUser.Username)
	}
	if updatedUser.DisplayName != nil && *updatedUser.DisplayName != "" {
		currentUser.DisplayName = *updatedUser.DisplayName
	}
	if updatedUser.Email != nil && *updatedUser.Email != "" {
		currentUser.Email = *updatedUser.Email
	}
	if updatedUser.AvatarURL != nil && *updatedUser.AvatarURL != "" {
		currentUser.AvatarURL = *updatedUser.AvatarURL
	}
	if updatedUser.PublicKey != nil && *updatedUser.PublicKey != "" {
		currentUser.PublicKey = *updatedUser.PublicKey
	}
	if updatedUser.PrivateKey != nil && *updatedUser.PrivateKey != "" {
		currentUser.PrivateKey = *updatedUser.PrivateKey
	}

	// update user
	err = uh.userService.UpdateUser(r.Context(), currentUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UpdateUserResponse{
		ID:          currentUser.ID,
		Username:    currentUser.Username,
		DisplayName: currentUser.DisplayName,
		Email:       currentUser.Email,
		AvatarURL:   currentUser.AvatarURL,
		ActorID:     currentUser.ActorID,
		PublicKey:   currentUser.PublicKey,
		CreatedAt:   currentUser.CreatedAt,
		UpdatedAt:   currentUser.UpdatedAt,
	})
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// get user id from request URL params
	idFromParam := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idFromParam)
	if err != nil {
		http.Error(w, "invalid user ID format from request parameter", http.StatusBadRequest)
		return
	}

	// get authenticated user from context (set by JWT middleware)
	userIDFromToken, err := uh.authHandler.GetUserIDByAuthTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if idFromParam != userIDFromToken {
		http.Error(w, "user can only delete their own profile", http.StatusForbidden)
		return
	}

	// delete user
	err = uh.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// response
	w.WriteHeader(http.StatusNoContent)
}
