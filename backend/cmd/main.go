package main

import (
	"sentimenta/internal/config"
	"sentimenta/internal/db"
	"sentimenta/internal/handlers"
	"sentimenta/internal/moodService"
	"sentimenta/internal/userService"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	preLogger, _ := zap.NewProduction()
	defer preLogger.Sync() // flushes buffer, if any
	logger := preLogger.Sugar()

	cfg := config.NewConfig()
	db := db.InitDB(cfg, logger)

	userRepo := userService.NewRepository(db)
	userService := userService.NewService(userRepo)

	moodRepo := moodService.NewRepository(db)
	moodService := moodService.NewService(moodRepo)

	userHandler := handlers.NewUserHandler(userService, logger)
	authHandler := handlers.NewAuthHandler(userService, *cfg, logger)
	moodHandler := handlers.NewMoodHandler(moodService, *cfg, logger)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/api/auth/login", authHandler.Login)
	e.POST("/api/auth/register", authHandler.Register)

	e.GET("/api/user/get", userHandler.GetUser)
	e.PATCH("/api/user/update", userHandler.PatchUpdateUser)
	e.PUT("/api/user/update/password", userHandler.PutUpdatePasswordUser)

	e.POST("/api/moods/add", moodHandler.PostAddMood)
	e.GET("/api/moods/get", moodHandler.GetMoods)
	e.PUT("/api/moods/update", moodHandler.PutUpdateMood)

	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}
