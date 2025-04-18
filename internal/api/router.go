package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"je-suis-ici-activitypub/internal/activitypub"
	"je-suis-ici-activitypub/internal/api/handlers"
	"je-suis-ici-activitypub/internal/api/middlewares"
	"je-suis-ici-activitypub/internal/services"
	"net/http"
)

func NewRouter(
	userService services.UserService,
	checkinService services.CheckinService,
	mediaService services.MediaService,
	apServerService *activitypub.ActivityPubServerService,
	tokenAuth *jwtauth.JWTAuth,
	serverHost string,
) http.Handler {
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	//r.Use(middleware.Recoverer)
	//r.Use(middleware.Timeout(60))

	// CORS setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"}, // allow browser read in response head
		AllowCredentials: true,
		MaxAge:           300, // 300sec
	}))

	// handlers
	authHandler := handlers.NewAuthHandler(userService, tokenAuth, serverHost)
	checkinHandler := handlers.NewCheckinHandler(userService, checkinService, mediaService, apServerService, serverHost)
	feedHandler := handlers.NewFeedHandler(checkinService)

	// public routes (no need JWT token)
	r.Group(func(r chi.Router) {
		// auth routes
		r.Route("/auth", func(r chi.Router) {
			authHandler.RegisterRoutes(r)
		})

		// ActivityPub routes
		r.Route("/.well-known", func(r chi.Router) {
			r.Get("/webfinger", nil) // implement WebFinger for ActivityPub
			r.Get("/nodeinfo", nil)  // implement NodeInfo for ActivityPub
		})
	})

	// routes with "/api" prefix
	r.Route("/api", func(r chi.Router) {
		// public routes (no need JWT token)
		r.Group(func(r chi.Router) {
			feedHandler.RegisterFeedRouters(r)
		})

		// protected routes (need JWT token)
		r.Group(func(r chi.Router) {
			// auth JWT middleware
			r.Use(middlewares.AuthJWT(tokenAuth))

			checkinHandler.RegisterCheckinRoutes(r)

			// activityPub user interaction routes
			r.Get("/users/{username}/activitypub-info", checkinHandler.GetUserActivityPubInfo)
			r.Post("/users/{sender_username}/send-checkin", checkinHandler.SendCheckinToUser)
			r.Get("/users/{username}/inbox", checkinHandler.GetUserInbox)
		})
	})

	return r
}
