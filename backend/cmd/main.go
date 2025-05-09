package main

import (
	"fmt"
	"sentimenta/internal/adviceService"
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
	jwt := security.NewJWT(cfg)
	oauth := auth.NewOAuth(cfg)

	userRepo := userService.NewRepository(db)
	userService := userService.NewService(userRepo)

	moodRepo := moodService.NewRepository(db)
	moodService := moodService.NewService(moodRepo)

	adviceRepo := adviceService.NewRepository(db)
	adviceServ := adviceService.NewService(adviceRepo, moodRepo, userRepo, cfg)

	userHandler := handlers.NewUserHandler(userService, cfg, logger)
	authHandler := handlers.NewAuthHandler(userService, cfg, logger, oauth, jwt)
	moodHandler := handlers.NewMoodHandler(moodService, cfg, logger)
	adviceHandler := handlers.NewAdviceHandler(adviceServ, logger)

	scheduler := adviceService.NewScheduler(adviceServ, logger)
	scheduler.Start()

	adviceServ.GenerateAdviceForAllUsers()

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
