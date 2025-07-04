package main

import (
	"fmt"
	"net/http"
	"sentimenta/internal/auth"
	"sentimenta/internal/config"
	"sentimenta/internal/db"
	"sentimenta/internal/handlers"
	"sentimenta/internal/metrics"
	middlewares "sentimenta/internal/middleware"
	"sentimenta/internal/repository"
	"sentimenta/internal/security"
	"sentimenta/internal/service"
	"sentimenta/internal/ws"

	_ "sentimenta/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

//	@title			Sentimenta API
//	@version		1.0
//	@description	API for mood tracker which uses AI

//	@host		localhost:8000
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth
//	@in							header
//	@name						Authorization

func main() {
	preLogger, _ := zap.NewDevelopment()
	defer func() {
		if err := preLogger.Sync(); err != nil {
			fmt.Printf("Failed to sync preLogger: %v\n", err)
		}
	}()
	logger := preLogger.Sugar()

	cfg := config.NewConfig()
	prometheusController := metrics.NewPrometheus()
	db := db.InitDB(cfg, logger, prometheusController)
	jwt := security.NewJWT(cfg)
	oauth := auth.NewOAuth(cfg)
	responser := handlers.NewResponser(prometheusController, logger)
	wsConnManager := ws.NewConnectionManager()

	userRepo := repository.NewUserRepository(db)
	moodRepo := repository.NewMoodRepository(db)
	adviceRepo := repository.NewAdviceRepository(db)

	userService := service.NewUserService(userRepo)
	adviceService := service.NewAdviceService(adviceRepo, moodRepo, userRepo, cfg, logger)
	moodService := service.NewMoodService(moodRepo, userRepo, adviceRepo, adviceService, logger, wsConnManager)

	wsHandler := handlers.NewWSHandler(logger, wsConnManager)
	userHandler := handlers.NewUserHandler(userService, cfg, logger, responser)
	authHandler := handlers.NewAuthHandler(userService, cfg, logger, oauth, jwt, responser)
	moodHandler := handlers.NewMoodHandler(moodService, cfg, logger, responser)
	adviceHandler := handlers.NewAdviceHandler(adviceService, logger, responser)
	statusHandler := handlers.NewStatusHandler()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.ALLOWED_ORIGINS,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	e.Use(middleware.Logger())
	e.Use(prometheusController.Middleware())

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

	e.GET("/ws", wsHandler.HandleWS, middlewares.NewJWTMiddleware(cfg, jwt))
	e.GET("/api/advice", adviceHandler.GetAdvice, middlewares.NewJWTMiddleware(cfg, jwt))
	e.GET("/api/status", statusHandler.GetStatus)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/swagger/*any", swagger.WrapHandler)

	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}
