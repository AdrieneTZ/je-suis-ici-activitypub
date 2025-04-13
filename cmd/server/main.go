package main

import (
	"context"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"je-suis-ici-activitypub/internal/api"
	"je-suis-ici-activitypub/internal/config"
	"je-suis-ici-activitypub/internal/db"
	"je-suis-ici-activitypub/internal/db/models"
	"je-suis-ici-activitypub/internal/services"
	"je-suis-ici-activitypub/internal/services/activitypub"
	"je-suis-ici-activitypub/internal/services/user"
	"je-suis-ici-activitypub/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("fail to load config: %w", err)
	}

	// init database connection
	database, err := db.NewDatabase(db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("fail to connect database: %w", err)
	}
	defer database.Close()

	fmt.Println("FINALLY success connect to database!!!")

	// execute database migrations
	err = db.ExecuteMigrations(db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("fail to execute database migrations: %w", err)
	}

	fmt.Println("success execute database migrations!!!")

	// init storage service (MinIO)
	storageService, err := storage.NewMinioServiceImplement(storage.MinioConfig{
		Endpoint:  cfg.MinioConfig.Endpoint,
		AccessKey: cfg.MinioConfig.AccessKey,
		SecretKey: cfg.MinioConfig.SecretKey,
		Bucket:    cfg.MinioConfig.Bucket,
		UseSSL:    cfg.MinioConfig.UseSSL,
	})
	if err != nil {
		log.Fatalf("fail to initialize storage service: %w", err)
	}

	fmt.Println("success initialize storage service!!!!")

	// init repositories
	userRepo := models.NewUserRepository(database.Pool)
	checkinRepo := models.NewCheckinRepository(database.Pool)
	mediaRepo := models.NewMediaRepository(database.Pool)

	// init services
	actorService := activitypub.NewActorService(userRepo)
	userService := user.NewUserService(userRepo, actorService)
	checkinService := services.NewCheckinService(checkinRepo, mediaRepo, storageService)
	mediaService := services.NewMediaService(mediaRepo, storageService)

	// init JWT auth
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWT.Secret), nil)

	// init router
	router := api.NewRouter(
		userService,
		checkinService,
		mediaService,
		tokenAuth,
		cfg.Server.Host,
	)

	// create HTTP server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           router,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    0,
	}

	// use channel to get operation signal
	signalChan := make(chan os.Signal, 1)
	// get specific signal
	// os.Interrupt: interrupt by Ctrl+C, syscall.SIGINT: interrupt by syscall, syscall.SIGTERM: request to terminate server
	// any signal above will be sent to signalChan
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start app server
	go func() {
		fmt.Println("start server")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("fail to start server: %w", err)
		}
	}()

	// get signal to shut down server
	<-signalChan
	fmt.Println("server is shutting down")

	// setup timeout to control shutting down server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// shut down server
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("fail to shut down server: %w", err)
	}

	fmt.Println("server is shut down!")
}
