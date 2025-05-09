package adviceService

import (
	"time"

	"gorm.io/gorm"
)

type AdviceRepository interface {
	GetAdvices(userID string) ([]Advice, error)
	GetAdvice(userID string, date time.Time) (Advice, error)
	CreateAdvice(a *Advice) error
	GetLastAdvice(userID string) (Advice, error)
}

type adviceRepository struct {
	db *gorm.DB
}

func (r *adviceRepository) CreateAdvice(advice *Advice) error {
	return r.db.Create(&advice).Error
}

func (r *adviceRepository) GetAdvices(userID string) ([]Advice, error) {
	var advices []Advice
	err := r.db.Find(&advices, "user_id = ?", userID).Error
	return advices, err
}

func (r *adviceRepository) GetAdvice(userID string, date time.Time) (Advice, error) {
	var advice Advice
	err := r.db.
		Where("DATE(date) = ?", date.Format("2006-01-02")).
		First(&advice).
		Error
	return advice, err
}

func (r *adviceRepository) GetLastAdvice(userID string) (Advice, error) {
	var advice Advice
	err := r.db.
		Where("user_id = ?", userID).
		Order("date DESC").
		First(&advice).
		Error

	return advice, err
}

func NewRepository(db *gorm.DB) AdviceRepository {
	return &adviceRepository{db: db}
}
