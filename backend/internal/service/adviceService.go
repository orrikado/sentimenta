package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sentimenta/internal/config"
	"sentimenta/internal/models"
	repo "sentimenta/internal/repository"
	"sentimenta/internal/utils"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdviceService interface {
	GetAdvice(userID string, date time.Time) (models.Advice, error)
	GetAdvices(userID string) ([]models.Advice, error)
	CreateAdvice(userID string, text string, date time.Time) (models.Advice, error)
	GetLastAdvice(userID string) (models.Advice, error)
	GenerateAdviceForAllUsers() error
	GenerateAdvice(userID int, date time.Time) (models.Advice, error)
}

type adviceService struct {
	repo     repo.AdviceRepository
	moodRepo repo.MoodRepository
	userRepo repo.UserRepository
	logger   *zap.SugaredLogger
	config   *config.Config
}

func (s *adviceService) CreateAdvice(userID string, text string, date time.Time) (models.Advice, error) {
	uidInt, err := strconv.Atoi(userID)
	if err != nil {
		return models.Advice{}, err
	}

	newAdvice := models.Advice{
		UserID: uidInt,
		Text:   text,
		Date:   date,
	}

	if err := s.repo.CreateAdvice(&newAdvice); err != nil {
		return models.Advice{}, err
	}

	return newAdvice, nil
}

func (s *adviceService) GetAdvices(userID string) ([]models.Advice, error) {
	return s.repo.GetAdvices(userID)
}

func (s *adviceService) GetAdvice(userID string, date time.Time) (models.Advice, error) {
	return s.repo.GetAdvice(userID, date)
}

func (s *adviceService) GenerateAdviceForAllUsers() error {
	users, err := s.userRepo.GetAllUsers() // нужен userRepo
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.IsActive {
			uidStr := fmt.Sprintf("%v", user.Uid)
			if _, err := s.repo.GetAdvice(uidStr, time.Now()); err == nil {
				continue
			}

			loc, err := time.LoadLocation(user.Timezone)
			if err != nil {
				s.logger.Errorf("не удалось загрузить часовой пояс: %v", err)
				continue
			}
			userTime := time.Now().In(loc)
			advice, err := s.GenerateAdvice(user.Uid, userTime)
			if err != nil {
				s.logger.Errorf("не удалось сгенерировать advice: %v", err)
				continue
			}
			if err := s.repo.CreateAdvice(&advice); err != nil {
				s.logger.Errorf("не удалось добавить advice: %v", err)
				continue
			}
			user.IsActive = false
			if err := s.userRepo.UpdateUser(user.Uid, user); err != nil {
				s.logger.Errorf("не удалось обновить пользователя: %v", err)
				continue
			}
		}
	}

	return nil
}

func (s *adviceService) GenerateAdvice(userID int, date time.Time) (models.Advice, error) {
	uidStr := fmt.Sprintf("%v", userID)
	lastMoods, err := s.moodRepo.GetLastMoods(uidStr, 30)
	if err != nil {
		return models.Advice{}, err
	}
	lastAdvice, err := s.repo.GetLastAdvice(uidStr)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		lastAdvice = models.Advice{Text: ""}
	}

	// Строим DTO
	var moodsDTO []models.MoodAdd
	for _, m := range lastMoods {
		moodsDTO = append(moodsDTO, models.MoodAdd{
			Score:       m.Score,
			Emotions:    m.Emotions,
			Description: m.Description,
			Date:        m.Date, // нужный формат
		})
	}

	payload := models.AdviceRequest{
		PreviousAdvice: lastAdvice.Text,
		Moods:          moodsDTO,
	}

	// Преобразуем в JSON (например, для отправки куда-то)
	// systemPrompt := `Ты — умный и заботливый помощник. На основе настроения пользователя за последние дни и предыдущего совета, составь короткий полезный совет по улучшению настроения. Учитывай, если настроение падало или росло.`

	userContentBytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return models.Advice{}, err
	}

	reqBody := utils.OpenRouterRequest{
		Model: "google/gemma-3-27b-it:free",
		Messages: []utils.OpenRouterMessage{
			{
				Role:    "system",
				Content: s.config.SYSTEM_PROMPT,
			},
			{
				Role:    "user",
				Content: string(userContentBytes),
			},
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return models.Advice{}, err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		return models.Advice{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.AI_API_KEY)
	// req.Header.Set("HTTP-Referer", "https://your-site.com") // если надо
	// req.Header.Set("X-Title", "MoodTrackerApp")             // если надо

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Advice{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Advice{}, err
	}

	// Можно логировать полный ответ
	fmt.Println("OpenRouter response:", string(bodyBytes))

	// Разбор результата
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return models.Advice{}, err
	}

	if len(result.Choices) == 0 {
		return models.Advice{}, fmt.Errorf("AI вернул пустой результат")
	}

	generatedText := result.Choices[0].Message.Content

	// Сохраняем результат как Advice
	advice := models.Advice{
		UserID: userID,
		Date:   date,
		Text:   generatedText,
	}
	if err != nil {
		return models.Advice{}, err
	}

	return advice, nil
}

func (s *adviceService) GetLastAdvice(userID string) (models.Advice, error) {
	return s.repo.GetLastAdvice(userID)
}
func NewAdviceService(repo repo.AdviceRepository, moodRepo repo.MoodRepository, userRepo repo.UserRepository, config *config.Config, logger *zap.SugaredLogger) AdviceService {
	return &adviceService{repo: repo, moodRepo: moodRepo, userRepo: userRepo, config: config, logger: logger}
}
