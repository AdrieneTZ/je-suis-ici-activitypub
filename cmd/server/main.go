package main

import (
	"context"
	"fmt"
	"je-suis-ici-activitypub/internal/activitypub"
	"je-suis-ici-activitypub/internal/api"
	"je-suis-ici-activitypub/internal/config"
	"je-suis-ici-activitypub/internal/db"
	"je-suis-ici-activitypub/internal/db/models"
	"je-suis-ici-activitypub/internal/services"
	"je-suis-ici-activitypub/internal/storage"
	"je-suis-ici-activitypub/internal/tracing"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

func main() {
	// init logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("fail to initialize logger: %w", err)
	}
	defer logger.Sync()

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("fail to load config", zap.Error(err))
	}

	// init jaeger tracer
	if cfg.Jaeger.Enable {
		tp, err := tracing.InitJaeger(&cfg.Jaeger)
		if err != nil {
			logger.Fatal("fail to init jaeger tracer: %w", zap.Error(err))
		}

		// tp is not nil means tracer is enabled
		if tp != nil {
			// when main function is closed, shutdown jaeger tracer provider
			defer func() {
				err := tp.Shutdown(context.Background())
				if err != nil {
					logger.Error("fail to shutdown jaeger tracer provider: %w", zap.Error(err))
				}
			}()
		}

		logger.Info("success init jaeger tracer!!!")
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
		logger.Fatal("fail to connect database: %w", zap.Error(err))
	}
	defer database.Close()

	logger.Info("success connect to database!!!")

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
		logger.Fatal("fail to execute database migrations: %w", zap.Error(err))
	}

	logger.Info("success execute database migrations!!!")

	// init storage service (MinIO)
	storageService, err := storage.NewMinioServiceImplement(storage.MinioConfig{
		Endpoint:  cfg.MinioConfig.Endpoint,
		AccessKey: cfg.MinioConfig.AccessKey,
		SecretKey: cfg.MinioConfig.SecretKey,
		Bucket:    cfg.MinioConfig.Bucket,
		UseSSL:    cfg.MinioConfig.UseSSL,
	})
	if err != nil {
		logger.Fatal("fail to initialize storage service: %w", zap.Error(err))
	}

	logger.Info("success initialize storage service!!!!")

	// init repositories
	userRepo := models.NewUserRepository(database.Pool)
	checkinRepo := models.NewCheckinRepository(database.Pool)
	mediaRepo := models.NewMediaRepository(database.Pool)
	activityRepo := activitypub.NewActivityPubRepository(database.Pool)
	followerRepo := activitypub.NewFollowerRepository(database.Pool)

	// init services
	actorService := activitypub.NewActorService(userRepo)
	userService := services.NewUserService(userRepo, actorService)
	checkinService := services.NewCheckinService(checkinRepo, mediaRepo, storageService)
	mediaService := services.NewMediaService(mediaRepo, storageService)

	// init ActivityPub services
	apClientService := activitypub.NewActivityPubClientService(nil)
	apServerService := activitypub.NewActivityPubServerService(
		activityRepo,
		followerRepo,
		userRepo,
		checkinRepo,
		actorService,
		apClientService,
		cfg.Server.Host,
	)

	// init JWT auth
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWT.Secret), nil)

	// init router
	router := api.NewRouter(
		logger,
		userService,
		checkinService,
		mediaService,
		apServerService,
		actorService,
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
			logger.Fatal("fail to start server: %w", zap.Error(err))
		}
	}()

	// get signal to shut down server
	<-signalChan
	logger.Info("server is shutting down")

	// setup timeout to control shutting down server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// shut down server
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Fatal("server force to shutdown", zap.Error(err))
	}

	logger.Info("server is shut down!")
}
