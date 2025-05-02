package moodservice

import "gorm.io/gorm"

type MoodRepository interface {
	GetMoods(user_id string) ([]Mood, error)
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
	return r.db.Delete(&Mood{}, "id = ?", id).Error
}

func (r *moodRepository) GetMoods(user_id string) ([]Mood, error) {
	var moods []Mood
	err := r.db.Find(&moods).Error
	return moods, err
}

func (r *moodRepository) UpdateMood(mood *Mood) error {
	return r.db.Save(&mood).Error
}

func NewRepository(db *gorm.DB) MoodRepository {
	return &moodRepository{db: db}
}
