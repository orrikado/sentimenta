package moodService

import (
	"strconv"
)

type MoodService interface {
	GetMoods(userID string) ([]Mood, error)
	CreateMood(userID string, score int, emotions string, description string) (Mood, error)
	UpdateMood(m *Mood) error
	DeleteMood(id string) error
}

type moodService struct {
	repo MoodRepository
}

func (s *moodService) CreateMood(userID string, score int, emotions string, description string) (Mood, error) {
	uidInt, err := strconv.Atoi(userID)
	if err != nil {
		return Mood{}, err
	}

	newMood := Mood{
		Score:       score,
		Emotions:    emotions,
		Description: description,
		UserId:      uidInt,
	}

	if err := s.repo.CreateMood(&newMood); err != nil {
		return Mood{}, err
	}
	return newMood, nil
}

func (s *moodService) DeleteMood(id string) error {
	return s.repo.DeleteMood(id)
}

func (s *moodService) GetMoods(userID string) ([]Mood, error) {
	return s.repo.GetMoods(userID)
}

func (s *moodService) UpdateMood(m *Mood) error {
	return s.repo.UpdateMood(m)
}

func NewService(repo MoodRepository) MoodService {
	return &moodService{repo: repo}
}
