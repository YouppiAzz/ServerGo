package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"server/config"
	"server/handlers"
	"server/middleware"
	"server/server"
)

func main() {
	cfg := config.Load()
	noDB := os.Getenv("NO_DB") == "true"

	var (
		authHandler   *handlers.AuthHandler
		userHandler   *handlers.UserHandler
		healthHandler *handlers.HealthHandler
	)

	if !noDB {
		// You would initialize your DB here and pass it to handlers
		// dbConn, err := database.NewDB(cfg.DatabaseURL)
		// ...
		// For this refactor, we'll skip DB setup for simplicity
		authHandler = handlers.NewAuthHandler(nil, cfg.JWTSecret)
		userHandler = handlers.NewUserHandler(nil)
	} else {
		authHandler = handlers.NewAuthHandler(nil, cfg.JWTSecret)
		userHandler = handlers.NewUserHandler(nil)
	}
	healthHandler = handlers.NewHealthHandler()

	srv := server.NewServer(cfg.Port)

	srv.Use(middleware.CORS())
	srv.Use(middleware.Logger())
	srv.Use(middleware.Security())
	srv.Use(middleware.RateLimiter(100))

	srv.GET("/health", healthHandler.Health)
	srv.POST("/auth/register", authHandler.Register)
	srv.POST("/auth/login", authHandler.Login)
	srv.GET("/auth/me", middleware.RequireAuth(cfg.JWTSecret)(userHandler.GetProfile))
	srv.PUT("/auth/me", middleware.RequireAuth(cfg.JWTSecret)(userHandler.UpdateProfile))
	srv.GET("/users", middleware.RequireAuth(cfg.JWTSecret)(userHandler.ListUsers))

	log.Printf("Starting server on port %s", cfg.Port)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
	srv.Stop()
}
