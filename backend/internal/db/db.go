package db

import (
	"fmt"
	c "sentimenta/internal/config"
	"sentimenta/internal/metrics"
	"sentimenta/internal/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitDB(cfg *c.Config, log *zap.SugaredLogger, prometheus *metrics.Prometheus) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v",
		cfg.POSTGRES_HOST, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_DB, cfg.POSTGRES_PORT)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &metrics.GormLogger{
			Prometheus: prometheus,
			Interface:  gormLogger.Default.LogMode(gormLogger.Info),
		},
	})
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	log.Info("БД: Подключение | Успешно.")
	if err := db.AutoMigrate(models.User{}, models.Mood{}, models.Advice{}); err != nil {
		log.Fatalf("Не удалось произвести миграцию: %v", err)
	}
	log.Info("БД: Автомиграция | Успешно.")
	return db
}
