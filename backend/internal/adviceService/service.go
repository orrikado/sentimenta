package adviceService

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sentimenta/internal/config"
	ms "sentimenta/internal/moodService"
	us "sentimenta/internal/userService"
	"sentimenta/internal/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AdviceService interface {
	GetAdvice(userID string, date time.Time) (Advice, error)
	GetAdvices(userID string) ([]Advice, error)
	CreateAdvice(userID string, text string, date time.Time) (Advice, error)
	GetLastAdvice(userID string) (Advice, error)
	GenerateAdviceForAllUsers() error
	GenerateAdvice(userID int, date time.Time) (Advice, error)
}

type adviceService struct {
	repo     AdviceRepository
	moodRepo ms.MoodRepository
	userRepo us.UserRepository
	config   *config.Config
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

func (s *adviceService) GenerateAdviceForAllUsers() error {
	users, err := s.userRepo.GetAllUsers() // нужен userRepo
	if err != nil {
		return err
	}
	for _, user := range users {
		advice, err := s.GenerateAdvice(user.Uid, time.Now())
		if err != nil {
			continue
		}
		err = s.repo.CreateAdvice(&advice)
		if err != nil {
			continue
		}
	}

	return nil
}

func (s *adviceService) GenerateAdvice(userID int, date time.Time) (Advice, error) {
	uidStr := fmt.Sprintf("%v", userID)
	lastMoods, err := s.moodRepo.GetLastMoods(uidStr, 30)
	if err != nil {
		return Advice{}, err
	}
	lastAdvice, err := s.repo.GetLastAdvice(uidStr)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		lastAdvice = Advice{Text: ""}
	}

	// Строим DTO
	var moodsDTO []ms.MoodAdd
	for _, m := range lastMoods {
		moodsDTO = append(moodsDTO, ms.MoodAdd{
			Score:       m.Score,
			Emotions:    m.Emotions,
			Description: m.Description,
			Date:        m.Date, // нужный формат
		})
	}

	payload := AdviceRequest{
		PreviousAdvice: lastAdvice.Text,
		Moods:          moodsDTO,
	}

	// Преобразуем в JSON (например, для отправки куда-то)
	systemPrompt := `Ты — умный и заботливый помощник. На основе настроения пользователя за последние дни и предыдущего совета, составь короткий полезный совет по улучшению настроения. Учитывай, если настроение падало или росло.`

	userContentBytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return Advice{}, err
	}

	reqBody := utils.OpenRouterRequest{
		Model: "openai/gpt-4o",
		Messages: []utils.OpenRouterMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: string(userContentBytes),
			},
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return Advice{}, err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		return Advice{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.AI_API_KEY)
	// req.Header.Set("HTTP-Referer", "https://your-site.com") // если надо
	// req.Header.Set("X-Title", "MoodTrackerApp")             // если надо

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Advice{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Advice{}, err
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
		return Advice{}, err
	}

	if len(result.Choices) == 0 {
		return Advice{}, fmt.Errorf("OpenRouter вернул пустой результат")
	}

	generatedText := result.Choices[0].Message.Content

	// Сохраняем результат как Advice
	advice := Advice{
		UserID: userID,
		Date:   date,
		Text:   generatedText,
	}
	err = s.repo.CreateAdvice(&advice)
	if err != nil {
		return Advice{}, err
	}

	return advice, nil
}

func (s *adviceService) GetLastAdvice(userID string) (Advice, error) {
	return s.repo.GetLastAdvice(userID)
}
func NewService(repo AdviceRepository, moodRepo ms.MoodRepository, userRepo us.UserRepository, config *config.Config) AdviceService {
	return &adviceService{repo: repo, moodRepo: moodRepo, userRepo: userRepo, config: config}
}
