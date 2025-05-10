package main

import (
	"fmt"
	"sentimenta/internal/auth"
	"sentimenta/internal/config"
	"sentimenta/internal/db"
	"sentimenta/internal/handlers"
	middlewares "sentimenta/internal/middleware"
	"sentimenta/internal/repository"
	"sentimenta/internal/scheduler"
	"sentimenta/internal/security"
	"sentimenta/internal/service"

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
	jwt := security.NewJWT(cfg)
	oauth := auth.NewOAuth(cfg)

	userRepo := repository.NewUserRepository(db)
	moodRepo := repository.NewMoodRepository(db)
	adviceRepo := repository.NewAdviceRepository(db)

	userService := service.NewUserService(userRepo)
	adviceService := service.NewAdviceService(adviceRepo, moodRepo, userRepo, cfg, logger)
	moodService := service.NewMoodService(moodRepo, userRepo, adviceRepo, adviceService, logger)

	userHandler := handlers.NewUserHandler(userService, cfg, logger)
	authHandler := handlers.NewAuthHandler(userService, cfg, logger, oauth, jwt)
	moodHandler := handlers.NewMoodHandler(moodService, cfg, logger)
	adviceHandler := handlers.NewAdviceHandler(adviceService, logger)

	scheduler := scheduler.NewScheduler(adviceService, logger)
	scheduler.Start()

	if err := adviceService.GenerateAdviceForAllUsers(); err != nil {
		logger.Errorf("Не удалось сгенерировать advice: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/api/auth/login", authHandler.Login)
	e.POST("/api/auth/register", authHandler.Register)

	e.POST("/api/auth/google/callback", authHandler.GoogleAuthCallback)
	e.POST("/api/auth/github/callback", authHandler.GithubAuthCallback)

	userGroup := e.Group("/api/user")
	userGroup.Use(middlewares.NewJWTMiddleware(cfg, jwt))
	userGroup.GET("/get", userHandler.GetUser)
	userGroup.PATCH("/update", userHandler.PatchUpdateUser)
	userGroup.PUT("/update/password", userHandler.PutUpdatePasswordUser)

	moodGroup := e.Group("/api/moods")
	moodGroup.Use(middlewares.NewJWTMiddleware(cfg, jwt))
	moodGroup.POST("/add", moodHandler.PostAddMood)
	moodGroup.GET("/get", moodHandler.GetMoods)
	moodGroup.PUT("/update", moodHandler.PutUpdateMood)

	e.GET("/api/advice", adviceHandler.GetAdvice, middlewares.NewJWTMiddleware(cfg, jwt))

	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}
