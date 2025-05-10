package service

import (
	m "sentimenta/internal/models"
	repo "sentimenta/internal/repository"
	"strconv"
	"time"
)

type MoodService interface {
	GetMoods(userID string) ([]m.Mood, error)
	CreateMood(userID string, score int16, emotions, description string, date time.Time) (m.Mood, error)
	UpdateMood(userID string, m *m.Mood) error
	DeleteMood(id string) error
}

type moodService struct {
	repo     repo.MoodRepository
	userRepo repo.UserRepository
}

func (s *moodService) CreateMood(userID string, score int16, emotions, description string, date time.Time) (m.Mood, error) {
	uidInt, err := strconv.Atoi(userID)
	if err != nil {
		return m.Mood{}, err
	}

	newMood := m.Mood{
		Score:       score,
		Emotions:    emotions,
		Description: description,
		UserId:      uidInt,
		Date:        date,
	}

	if err := s.repo.CreateMood(&newMood); err != nil {
		return m.Mood{}, err
	}

	dateStr := date.Format("2006-01-02")
	dateNowStr := time.Now().Format("2006-01-02")

	if dateStr != dateNowStr {
		userUpdate := m.User{
			Uid:      uidInt,
			IsActive: true,
		}
		if err := s.userRepo.UpdateUser(userUpdate); err != nil {
			return m.Mood{}, err
		}
	}
	return newMood, nil
}

func (s *moodService) DeleteMood(id string) error {
	return s.repo.DeleteMood(id)
}

func (s *moodService) GetMoods(userID string) ([]m.Mood, error) {
	return s.repo.GetMoods(userID)
}

func (s *moodService) UpdateMood(userID string, m *m.Mood) error {
	uidInt, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	m.UserId = uidInt
	return s.repo.UpdateMood(m)
}

func NewMoodService(repo repo.MoodRepository) MoodService {
	return &moodService{repo: repo}
}
