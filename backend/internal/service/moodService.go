package service

import (
	"encoding/json"
	m "sentimenta/internal/models"
	repo "sentimenta/internal/repository"
	"sentimenta/internal/ws"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type moodService struct {
	repo       repo.MoodRepository
	userRepo   repo.UserRepository
	adviceRepo repo.AdviceRepository
	adviceServ AdviceService
	connMgr    *ws.ConnectionManager
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

	if user.UseAI {
		loc, err := time.LoadLocation(user.Timezone)
		if err != nil {
			s.logger.Errorf("не удалось загрузить часовой пояс: %v", err)
		}
		dateStr := date.Format("2006-01-02")
		dateNowStr := time.Now().In(loc).Format("2006-01-02")
		dateYesterdayStr := time.Now().AddDate(0, 0, -1).In(loc).Format("2006-01-02")

		if dateStr == dateNowStr || dateStr == dateYesterdayStr {
			if err := s.userRepo.UpdateUser(uidInt, map[string]interface{}{"is_active": true}); err != nil {
				return m.Mood{}, err
			}
			go func() {
				advice, err := s.adviceServ.GenerateAdvice(uidInt, date)
				if err != nil {
					s.logger.Errorf("не удалось сгенерировать advice: %v", err)
					return
				}
				if err := s.adviceRepo.CreateAdvice(&advice); err != nil {
					s.logger.Errorf("не удалось добавить advice: %v", err)
				}

				adviceJson, err := json.Marshal(advice)
				if err != nil {
					s.logger.Errorf("не удался Marshal adviceJson: %v", err)
					return
				}

				if err := s.connMgr.Send(userID, string(adviceJson)); err != nil {
					s.logger.Errorf("не удалось отправить advice по WS: %v", err)
				}

			}()
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

func NewMoodService(
	repo repo.MoodRepository,
	userRepo repo.UserRepository,
	adviceRepo repo.AdviceRepository,
	adviceServ AdviceService,
	logger *zap.SugaredLogger,
	wsConnMgr *ws.ConnectionManager,
) *moodService {
	return &moodService{
		repo:       repo,
		userRepo:   userRepo,
		adviceRepo: adviceRepo,
		adviceServ: adviceServ,
		logger:     logger,
		connMgr:    wsConnMgr,
	}
}
