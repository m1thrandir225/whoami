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
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/handlers"
	"github.com/m1thrandir225/whoami/internal/repositories"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/services"
	"github.com/m1thrandir225/whoami/internal/util"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer connPool.Close()

	rateLimiter, err := security.NewRateLimiter(config.RedisURL)
	if err != nil {
		log.Fatalf("Could not create rate limiter: %v", err)
	}
	defer rateLimiter.Close()

	dbStore := db.NewStore(connPool)
	tokenMaker, err := security.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Could not create tokenMaker: %v", err)
	}

	/*
	* Users
	 */
	userRepository := repositories.NewUserRepository(dbStore)
	userService := services.NewUserService(userRepository)
	handler := handlers.NewHTTPHandler(userService, tokenMaker, rateLimiter, config)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers.SetupRoutes(router, handler)

	httpServer := &http.Server{
		Addr:    config.HTTPServerAddress,
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	cancel()

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
}
