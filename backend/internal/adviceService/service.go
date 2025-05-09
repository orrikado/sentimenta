package adviceService

import (
	"strconv"
	"time"
)

type AdviceService interface {
	GetAdvice(userID string, date time.Time) (Advice, error)
	GetAdvices(userID string) ([]Advice, error)
	CreateAdvice(userID string, text string, date time.Time) (Advice, error)
}

type adviceService struct {
	repo AdviceRepository
}

func (s *adviceService) CreateAdvice(userID string, text string, date time.Time) (Advice, error) {
	uidInt, err := strconv.Atoi(userID)
	if err != nil {
		return Advice{}, err
	}

	newAdvice := Advice{
		UserID: uidInt,
		Text:   text,
		Date:   date,
	}

	if err := s.repo.CreateAdvice(&newAdvice); err != nil {
		return Advice{}, err
	}

	return newAdvice, nil
}

func (s *adviceService) GetAdvices(userID string) ([]Advice, error) {
	return s.repo.GetAdvices(userID)
}

func (s *adviceService) GetAdvice(userID string, date time.Time) (Advice, error) {
	return s.repo.GetAdvice(userID, date)
}

func NewService(repo AdviceRepository) AdviceService {
	return &adviceService{repo: repo}
}
