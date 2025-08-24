package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/whoami/cmd/pkg/redis"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/handlers"
	"github.com/m1thrandir225/whoami/internal/repositories"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/services"
	"github.com/m1thrandir225/whoami/internal/util"
)

func main() {
	/**
	* Load config
	 */
	ctx, cancel := context.WithCancel(context.Background())

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	/**
	* Connect to database
	 */
	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer connPool.Close()

	/**
	* Connect to Redis
	 */
	redisClient, err := redis.NewRedisClient(config.RedisURL)
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	defer redisClient.Close()

	/**
	* Create rate limiter
	 */
	rateLimiter, err := security.NewRateLimiter(redisClient)
	if err != nil {
		log.Fatalf("Could not create rate limiter: %v", err)
	}
	defer rateLimiter.Close()

	/**
	* Create token maker
	 */
	tokenMaker, err := security.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Could not create tokenMaker: %v", err)
	}

	/**
	* Create database store
	 */
	dbStore := db.NewStore(connPool)

	/*
	* Repositories
	 */
	userRepository := repositories.NewUserRepository(dbStore)
	accountLockoutRepository := repositories.NewAccountLockoutRepository(dbStore)
	suspiciousActivityRepository := repositories.NewSuspiciousActivityRepository(dbStore)
	passwordHistoryRepository := repositories.NewPasswordHistoryRepository(dbStore)
	emailVerificationRepository := repositories.NewEmailVerificationRepository(dbStore)
	passwordResetRepository := repositories.NewPasswordResetRepository(dbStore)
	loginAttemptsRepository := repositories.NewLoginAttemptsRepository(dbStore)

	/*
	* Services
	 */
	userService := services.NewUserService(userRepository)
	securityService := services.NewSecurityService(
		loginAttemptsRepository,
		suspiciousActivityRepository,
		accountLockoutRepository,
		userRepository,
	)
	passwordSecurityService := services.NewPasswordSecurityService(
		passwordHistoryRepository,
		userRepository,
	)
	emailService := services.NewEmailService(
		emailVerificationRepository,
		userRepository,
		&config,
	)
	passwordResetService := services.NewPasswordResetService(
		passwordResetRepository,
		userRepository,
		passwordSecurityService,
		&config,
	)
	tokenBlacklist := security.NewTokenBlacklist(redisClient)
	sessionService := services.NewSessionService(redisClient, tokenBlacklist)

	/**
	* Create HTTP handler
	 */
	handler := handlers.NewHTTPHandler(
		userService,
		securityService,
		passwordSecurityService,

		passwordResetService,
		emailService,
		tokenMaker,
		tokenBlacklist,
		sessionService,
		rateLimiter,
		config,
	)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers.SetupRoutes(router, handler)

	httpServer := &http.Server{
		Addr:    config.HTTPServerAddress,
		Handler: router,
	}

	/**
	* Start HTTP server
	 */
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	/**
	* Wait for shutdown signal
	 */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	cancel()

	/**
	* Shutdown HTTP server
	 */
	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
}
