package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"je-suis-ici-activitypub/internal/api/handlers"
	"je-suis-ici-activitypub/internal/services/user"
	"net/http"
)

func NewRouter(
	userService user.UserService,
	tokenAuth *jwtauth.JWTAuth,
	serverHost string,
) http.Handler {
	r := chi.NewRouter()

	// CORS setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"}, // allow browser read in response head
		AllowCredentials: true,
		MaxAge:           300, // 300sec
	}))

	// handlers
	authHandler := handlers.NewAuthHandler(userService, tokenAuth, serverHost)

	// open route (no need JWT token)
	r.Group(func(r chi.Router) {
		// auth routes
		r.Route("/auth", func(r chi.Router) {
			authHandler.RegisterRoutes(r)
		})
	})

	// protected routes (need JWT token)
	r.Group(func(r chi.Router) {
	})

	return r
}
