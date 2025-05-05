package main

import (
	"fmt"
	"sentimenta/internal/auth"
	"sentimenta/internal/config"
	"sentimenta/internal/db"
	"sentimenta/internal/handlers"
	middlewares "sentimenta/internal/middleware"
	"sentimenta/internal/moodService"
	"sentimenta/internal/security"
	"sentimenta/internal/userService"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	preLogger, _ := zap.NewDevelopment()
	defer func() {
		if err := preLogger.Sync(); err != nil {
			fmt.Printf("Failed to sync preLogger: %v\n", err)
		}
	}()
	logger := preLogger.Sugar()

	cfg := config.NewConfig()
	db := db.InitDB(cfg, logger)
	jwt := security.NewJWT(*cfg)
	oauthConfig := auth.NewOAuthConfig(*cfg)

	userRepo := userService.NewRepository(db)
	userService := userService.NewService(userRepo)

	moodRepo := moodService.NewRepository(db)
	moodService := moodService.NewService(moodRepo)

	userHandler := handlers.NewUserHandler(userService, logger)
	authHandler := handlers.NewAuthHandler(userService, *cfg, logger, oauthConfig, jwt)
	moodHandler := handlers.NewMoodHandler(moodService, *cfg, logger)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/api/auth/login", authHandler.Login)
	e.POST("/api/auth/register", authHandler.Register)

	e.POST("/api/auth/google/callback", authHandler.GoogleAuthCallback)

	userGroup := e.Group("/api/user")
	userGroup.Use(middlewares.NewJWTMiddleware(*cfg, jwt))
	userGroup.GET("/get", userHandler.GetUser)
	userGroup.PATCH("/update", userHandler.PatchUpdateUser)
	userGroup.PUT("/update/password", userHandler.PutUpdatePasswordUser)

	moodGroup := e.Group("/api/moods")
	moodGroup.Use(middlewares.NewJWTMiddleware(*cfg, jwt))
	moodGroup.POST("/add", moodHandler.PostAddMood)
	moodGroup.GET("/get", moodHandler.GetMoods)
	moodGroup.PUT("/update", moodHandler.PutUpdateMood)

	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}
