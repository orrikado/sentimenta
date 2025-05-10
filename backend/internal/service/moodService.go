package service

import (
	"fmt"
	m "sentimenta/internal/models"
	repo "sentimenta/internal/repository"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type MoodService interface {
	GetMoods(userID string) ([]m.Mood, error)
	CreateMood(userID string, score int16, emotions, description string, date time.Time) (m.Mood, error)
	UpdateMood(userID string, m *m.Mood) error
	DeleteMood(id string) error
}

type moodService struct {
	repo       repo.MoodRepository
	userRepo   repo.UserRepository
	adviceRepo repo.AdviceRepository
	adviceServ AdviceService
	logger     *zap.SugaredLogger
}

func (s *moodService) CreateMood(userID string, score int16, emotions, description string, date time.Time) (m.Mood, error) {
	uidInt, err := strconv.Atoi(userID)
	if err != nil {
		return m.Mood{}, err
	}
	user, err := s.userRepo.GetUser(userID)
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
	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		s.logger.Errorf("не удалось загрузить часовой пояс: %v", err)
	}
	dateNowStr := time.Now().In(loc).Format("2006-01-02")

	fmt.Println(dateStr, dateNowStr)
	if dateStr == dateNowStr {
		if err := s.userRepo.UpdateUser(uidInt, map[string]interface{}{"is_active": true}); err != nil {
			return m.Mood{}, err
		}
		go func() {
			advice, err := s.adviceServ.GenerateAdvice(uidInt, date)
			if err != nil {
				s.logger.Errorf("не удалось сгенерировать advice: %v", err)
			}
			if err := s.adviceRepo.CreateAdvice(&advice); err != nil {
				s.logger.Errorf("не удалось добавить advice: %v", err)
			}
		}()
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

func NewMoodService(
	repo repo.MoodRepository,
	userRepo repo.UserRepository,
	adviceRepo repo.AdviceRepository,
	adviceServ AdviceService,
	logger *zap.SugaredLogger,
) *moodService {
	return &moodService{
		repo:       repo,
		userRepo:   userRepo,
		adviceRepo: adviceRepo,
		adviceServ: adviceServ,
		logger:     logger,
	}
}
