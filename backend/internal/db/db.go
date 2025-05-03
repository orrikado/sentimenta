package db

import (
	"fmt"

	c "sentimenta/internal/config"
	"sentimenta/internal/moodService"
	"sentimenta/internal/userService"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(c *c.Config, logger *zap.SugaredLogger) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", c.POSTGRES_HOST, c.POSTGRES_USER, c.POSTGRES_PASSWORD, c.POSTGRES_DB, c.POSTGRES_PORT)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		logger.Fatalf("Не удалось подключится к БД: %v", err)
	}

	if err := db.AutoMigrate(userService.User{}, moodService.Mood{}); err != nil {
		logger.Fatalf("Не удалось произвести миграцию: %v", err)
	}

	return db
}
