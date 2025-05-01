package db

import (
	"fmt"
	"log"

	c "sentimenta/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(c *c.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", c.POSTGRES_HOST, c.POSTGRES_USER, c.POSTGRES_PASSWORD, c.POSTGRES_DB, c.POSTGRES_PORT)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("Не удалось подключится к БД: %v", err)
	}

	if err := db.AutoMigrate(); err != nil { // ! Добавить модели для миграции
		log.Fatalf("Не удалось произвести миграцию: %v", err)
	}

	return db
}
