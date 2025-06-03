package handlers

import (
	"encoding/json"
	"fmt"
	"je-suis-ici-activitypub/internal/db/models"
	"je-suis-ici-activitypub/internal/services"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("api/handlers/auth")

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
	// add a new root tracer and span
	ctx, span := tracer.Start(r.Context(), "AuthHandler.Register") // operation name
	defer span.End()                                               // close the span

	// record request attributes
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
	)

	// parse request body
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusBadRequest),
			attribute.String("error.type", "request_body_decode_error"),
			attribute.String("error.message", err.Error()),
		)
		span.SetStatus(codes.Error, "decode request body failed")

		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// valid request
	if req.Username == "" || req.Email == "" || req.Password == "" {
		span.RecordError(fmt.Errorf("missing required fields"))
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusBadRequest),
			attribute.String("error.type", "validation_error"),
			attribute.String("error.details", "missing_required_fields"),
			attribute.Bool("validation.username_empty", req.Username == ""),
			attribute.Bool("validation.email_empty", req.Email == ""),
			attribute.Bool("validation.password_empty", req.Password == ""),
		)
		span.SetStatus(codes.Error, "validation error")

		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// register a user account
	user, err := ah.userService.Register(ctx, ah.serverHost, req.Username, req.Email, req.Password)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusBadRequest),
			attribute.String("error.type", "registration_error"),
			attribute.String("error.message", err.Error()),
			attribute.String("req.username", req.Username),
			attribute.String("req.email", req.Email),
		)
		span.SetStatus(codes.Error, "registration failed")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// record user registration success
	span.SetAttributes(
		attribute.Int("response.status_code", http.StatusCreated),
		attribute.String("userID", user.ID.String()),
		attribute.String("username", user.Username),
	)
	span.SetStatus(codes.Ok, "user registered successfully")

	// TODO: refactor to a function
	// generate JWT token
	claims := map[string]interface{}{
		"user_id": user.ID.String(),                      // user uuid
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // expired time
	}
	_, tokenString, err := ah.tokenAuth.Encode(claims)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusInternalServerError),
			attribute.String("error.type", "jwt_token_generation_error"),
			attribute.String("error.message", err.Error()),
		)
		span.SetStatus(codes.Error, "jwt token generation failed")

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

// Login
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// add a new root tracer and span
	ctx, span := tracer.Start(r.Context(), "AuthHandler.Login") // operation name
	defer span.End()                                            // close the span

	// record request attributes
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
	)

	// parse login request
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusBadRequest),
			attribute.String("error.type", "request_body_decode_error"),
			attribute.String("error.message", err.Error()),
		)
		span.SetStatus(codes.Error, "decode request body failed")

		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// verify required request parameters
	if req.UsernameOrEmail == "" || req.Password == "" {
		span.RecordError(fmt.Errorf("missing required fields"))
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusBadRequest),
			attribute.String("error.type", "validation_error"),
			attribute.String("error.details", "missing_required_fields"),
			attribute.Bool("validation.username_or_email_empty", req.UsernameOrEmail == ""),
			attribute.Bool("validation.password_empty", req.Password == ""),
		)
		span.SetStatus(codes.Error, "validation error")

		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// verify user
	user, err := ah.userService.Authenticate(ctx, req.UsernameOrEmail, req.Password)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusBadRequest),
			attribute.String("error.type", "authentication_error"),
			attribute.String("error.message", err.Error()),
			attribute.String("req.username_or_email", req.UsernameOrEmail),
		)
		span.SetStatus(codes.Error, "authentication failed")

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
		span.RecordError(err)
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusInternalServerError),
			attribute.String("error.type", "generate_jwt_token_error"),
			attribute.String("error.message", err.Error()),
		)
		span.SetStatus(codes.Error, "jwt token generation failed")

		http.Error(w, "fail to generate JWT token", http.StatusInternalServerError)
		return
	}

	// record authentication success
	span.SetAttributes(
		attribute.Int("response.status_code", http.StatusOK),
		attribute.String("userID", user.ID.String()),
	)

	// response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: tokenString,
		User:  user,
	})
}

// GetUserIDByAuthTokenFromRequest
// JWT token is verified by the middleware, this function is used to get the claims from the context
// then get the user_id from the claims
func (ah *AuthHandler) GetUserIDByAuthTokenFromRequest(r *http.Request) (string, error) {
	// add a new root tracer and span
	ctx, span := tracer.Start(r.Context(), "AuthHandler.GetUserIDByAuthTokenFromRequest")
	defer span.End()

	// get JWT token claims from context
	// claims is like payload, and it's a map
	// claims contains the information (like user_id, exp...) in the JWT token
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusUnauthorized),
			attribute.String("error.type", "get_jwt_token_claims_error"),
			attribute.String("error.message", err.Error()),
		)
		span.SetStatus(codes.Error, "get JWT token claims failed")

		return "", fmt.Errorf("unauthorized: %w", err)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		span.RecordError(fmt.Errorf("invalid token"))
		// don't record claims content, avoid sensitive information leakage
		span.SetAttributes(
			attribute.String("error.type", "invalid_token"),
		)
		span.SetStatus(codes.Error, "invalid token")

		return "", fmt.Errorf("invalid token")
	}

	// when success, record userID
	span.SetAttributes(
		attribute.String("userID", userID),
	)

	return userID, nil
}
