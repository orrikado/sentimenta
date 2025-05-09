package adviceService

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Scheduler struct {
	cron    *cron.Cron
	service AdviceService
	logger  *zap.SugaredLogger
}

func NewScheduler(service AdviceService, logger *zap.SugaredLogger) *Scheduler {
	return &Scheduler{
		cron:    cron.New(),
		service: service,
		logger:  logger,
	}
}

func (s *Scheduler) Start() {
	// Каждый день в 3:00 ночи
	s.logger.Info("Scheduler started")
	_, err := s.cron.AddFunc("0 3 * * *", func() {
		err := s.service.GenerateAdviceForAllUsers()
		if err != nil {
			s.logger.Errorln("Ошибка при генерации советов:", err)
		}
	})
	if err != nil {
		s.logger.Errorln("Ошибка добавления cron-задачи:", err)
	}
	s.cron.Start()
}
