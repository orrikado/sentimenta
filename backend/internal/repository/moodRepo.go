package repository

import (
	m "sentimenta/internal/models"

	"gorm.io/gorm"
)

type MoodRepository interface {
	GetMoods(userID string) ([]m.Mood, error)
	GetLastMoods(userID string, limit int) ([]m.Mood, error)
	CreateMood(m *m.Mood) error
	UpdateMood(m *m.Mood) error
	DeleteMood(id string) error
}

type moodRepository struct {
	db *gorm.DB
}

func (r *moodRepository) CreateMood(mood *m.Mood) error {
	return r.db.Create(&mood).Error
}

func (r *moodRepository) DeleteMood(id string) error {
	return r.db.Delete(&m.Mood{}, "uid = ?", id).Error
}

func (r *moodRepository) GetMoods(userID string) ([]m.Mood, error) {
	var moods []m.Mood
	err := r.db.Find(&moods, "user_id = ?", userID).Error
	return moods, err
}

func (r *moodRepository) GetLastMoods(userID string, limit int) ([]m.Mood, error) {
	var moods []m.Mood
	err := r.db.Find(&moods, "user_id = ?", userID).Limit(limit).Order("date DESC").Error
	return moods, err
}

func (r *moodRepository) UpdateMood(mood *m.Mood) error {
	return r.db.Model(&m.Mood{}).Where("uid = ?", mood.Uid).Updates(mood).Error
}

func NewMoodRepository(db *gorm.DB) MoodRepository {
	return &moodRepository{db: db}
}
