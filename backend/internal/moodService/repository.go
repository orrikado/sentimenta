package moodService

import "gorm.io/gorm"

type MoodRepository interface {
	GetMoods(userID string) ([]Mood, error)
	GetLastMoods(userID string, limit int) ([]Mood, error)
	CreateMood(m *Mood) error
	UpdateMood(m *Mood) error
	DeleteMood(id string) error
}

type moodRepository struct {
	db *gorm.DB
}

func (r *moodRepository) CreateMood(mood *Mood) error {
	return r.db.Create(&mood).Error
}

func (r *moodRepository) DeleteMood(id string) error {
	return r.db.Delete(&Mood{}, "uid = ?", id).Error
}

func (r *moodRepository) GetMoods(userID string) ([]Mood, error) {
	var moods []Mood
	err := r.db.Find(&moods, "user_id = ?", userID).Error
	return moods, err
}

func (r *moodRepository) GetLastMoods(userID string, limit int) ([]Mood, error) {
	var moods []Mood
	err := r.db.Find(&moods, "user_id = ?", userID).Limit(limit).Order("date DESC").Error
	return moods, err
}

func (r *moodRepository) UpdateMood(mood *Mood) error {
	return r.db.Model(&Mood{}).Where("uid = ?", mood.Uid).Updates(mood).Error
}

func NewRepository(db *gorm.DB) MoodRepository {
	return &moodRepository{db: db}
}
