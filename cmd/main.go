package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/m1thrandir225/whoami/internal/db/sqlc"
	"github.com/m1thrandir225/whoami/internal/handlers"
	"github.com/m1thrandir225/whoami/internal/repositories"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/services"
	"github.com/m1thrandir225/whoami/internal/util"
	"log"
)

func main() {
	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	defer connPool.Close()

	dbStore := db.NewStore(connPool)
	tokenMaker, err := security.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Could not create tokenMaker: %v", err)
	}

	userRepository := repositories.NewUserRepository(dbStore)
	userService := services.NewUserService(userRepository)

	httpHandler := handlers.NewHTTPHandler(userService, tokenMaker, config)
}
