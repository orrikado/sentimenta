package repository

import (
	m "sentimenta/internal/models"
	"time"

	"gorm.io/gorm"
)

type adviceRepository struct {
	db *gorm.DB
}

func (r *adviceRepository) CreateAdvice(advice *m.Advice) error {
	return r.db.Create(&advice).Error
}

func (r *adviceRepository) GetAdvices(userID string) ([]m.Advice, error) {
	var advices []m.Advice
	err := r.db.Find(&advices, "user_id = ?", userID).Error
	return advices, err
}

func (r *adviceRepository) GetAdvice(userID string, date time.Time) (m.Advice, error) {
	var advice m.Advice
	err := r.db.
		Where("DATE(date) = ?", date.Format("2006-01-02")).
		First(&advice).
		Error
	return advice, err
}

func (r *adviceRepository) GetLastAdvice(userID string) (m.Advice, error) {
	var advice m.Advice
	err := r.db.
		Where("user_id = ?", userID).
		Order("date DESC").
		First(&advice).
		Error

	return advice, err
}

func NewAdviceRepository(db *gorm.DB) AdviceRepository {
	return &adviceRepository{db: db}
}
