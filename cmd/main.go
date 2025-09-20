package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/handlers"
	"github.com/m1thrandir225/whoami/internal/mail"
	"github.com/m1thrandir225/whoami/internal/oauth"
	"github.com/m1thrandir225/whoami/internal/repositories"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/services"
	"github.com/m1thrandir225/whoami/internal/util"
	"github.com/m1thrandir225/whoami/pkg/redis"
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

	/**
	* Create mail service
	 */
	mailService := mail.NewResendMailer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)

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
	auditLogsRepository := repositories.NewAuditLogsRepository(dbStore)
	userDevicesRepository := repositories.NewUserDevicesRepository(dbStore)
	dataExportsRepository := repositories.NewDataExportsRepository(dbStore)
	oauthAccountsRepository := repositories.NewOAuthAccountsRepository(dbStore)

	/*
	* OAuth Providers
	 */
	googleProvider := oauth.NewGoogleProvider(oauth.Config{
		ClientID:     config.GoogleOAuthClientID,
		ClientSecret: config.GoogleOAuthClientSecret,
		RedirectURL:  config.GoogleOAuthRedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
	})

	githubProvider := oauth.NewGitHubProvider(oauth.Config{
		ClientID:     config.GitHubOAuthClientID,
		ClientSecret: config.GitHubOAuthClientSecret,
		RedirectURL:  config.GitHubOAuthRedirectURL,
		Scopes:       []string{"read:user", "user:email"},
	})

	oauthProviders := handlers.OAuthProviders{
		Google: googleProvider,
		GitHub: githubProvider,
	}

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
		mailService,
	)
	passwordResetService := services.NewPasswordResetService(
		passwordResetRepository,
		userRepository,
		passwordSecurityService,
		mailService,
	)
	auditService := services.NewAuditService(auditLogsRepository)
	tokenBlacklist := security.NewTokenBlacklist(redisClient)
	sessionService := services.NewSessionService(redisClient, tokenBlacklist)
	userDevicesService := services.NewUserDevicesService(userDevicesRepository)

	exportDir := "./exports"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Fatalf("failed to create export directory: %v", err)
	}

	dataExportsService := services.NewDataExportsService(
		dataExportsRepository,
		userRepository,
		auditLogsRepository,
		loginAttemptsRepository,
		exportDir,
	)
	oauthService := services.NewOAuthService(
		oauthAccountsRepository,
		userRepository,
	)
	oauthTempService := services.NewOAuthTempService(redisClient)

	/**
	* Create HTTP handler
	 */
	handler := handlers.NewHTTPHandler(
		userService,
		securityService,
		passwordSecurityService,
		passwordResetService,
		emailService,
		auditService,
		userDevicesService,
		dataExportsService,
		oauthService,
		oauthProviders,
		tokenMaker,
		tokenBlacklist,
		sessionService,
		rateLimiter,
		oauthTempService,
		config,
	)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := dataExportsService.ProcessPendingExports(ctx); err != nil {
					log.Printf("failed to process pending exports: %v", err)
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := dataExportsService.CleanupExpiredExports(ctx); err != nil {
					log.Printf("failed to cleanup expired exports: %v", err)
				}
			}
		}
	}()

	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	//Cors Setup
	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	handlers.SetupRoutes(router, handler)

	httpAddress := fmt.Sprintf("%s:%d", config.HTTPServerAddress, config.HTTPPort)

	httpServer := &http.Server{
		Addr:    httpAddress,
		Handler: router,
	}

	/**
	* Start HTTP server
	 */
	go func() {
		if config.EnableTLS {
			if err := httpServer.ListenAndServeTLS(config.TLSCertFile, config.TLSKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("listen (TLS): %s\n", err)
			}
			return
		}
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
