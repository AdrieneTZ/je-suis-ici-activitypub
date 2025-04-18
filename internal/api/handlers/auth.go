package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"je-suis-ici-activitypub/internal/db/models"
	"je-suis-ici-activitypub/internal/services"
	"net/http"
	"time"
)

// AuthHandler handle auth requests
type AuthHandler struct {
	userService services.UserService
	tokenAuth   *jwtauth.JWTAuth
	serverHost  string
}

// NewAuthHandler
func NewAuthHandler(userService services.UserService, tokenAuth *jwtauth.JWTAuth, serverHost string) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		tokenAuth:   tokenAuth,
		serverHost:  serverHost,
	}
}

// RegisterRoutes register auth routes
func (ah *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", ah.Register)
	r.Post("/login", ah.Login)
}

// RegisterRequest
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

// AuthResponse
type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register user register an account
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// parse request body
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// valid request
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// register a user account
	user, err := ah.userService.Register(r.Context(), ah.serverHost, req.Username, req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// TODO: refactor to a function
	// generate JWT token
	claims := map[string]interface{}{
		"user_id": user.ID.String(),                      // user uuid
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // expired time
	}
	_, tokenString, err := ah.tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "fail to generate JWT token", http.StatusInternalServerError)
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		Token: tokenString,
		User:  user,
	})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// parse login request
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// verify required request parameters
	if req.UsernameOrEmail == "" || req.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// verify user
	user, err := ah.userService.Authenticate(r.Context(), req.UsernameOrEmail, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid credentials"})
		return
	}

	// TODO: refactor to a function
	// generate JWT token
	claims := map[string]interface{}{
		"user_id": user.ID.String(),                      // user uuid
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // expired time
	}
	_, tokenString, err := ah.tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "fail to generate JWT token", http.StatusInternalServerError)
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: tokenString,
		User:  user,
	})
}
