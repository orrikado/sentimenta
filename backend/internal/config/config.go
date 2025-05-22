package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string

	JWT_COOKIE_NAME string
	JWT_SECRET      string

	GOOGLE_CLIENT_ID       string
	GOOGLE_CLIENT_SECRET   string
	GOOGLE_CLIENT_CALLBACK string

	GITHUB_CLIENT_ID       string
	GITHUB_CLIENT_SECRET   string
	GITHUB_CLIENT_CALLBACK string

	SYSTEM_PROMPT string
	AI_API_KEY    string

	PASSWORD_LENGTH_MIN    int
	MOOD_DESC_LENGTH_MAX   int
	MOOD_EMOTES_LENGTH_MAX int

	JWT_HTTP_ONLY bool
	JWT_SECURE    bool
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("не удалось загрузить .env: %v\n", err)
	}

	passwordLenMin, err := strconv.Atoi(os.Getenv("PASSWORD_LENGTH_MIN"))
	if err != nil {
		fmt.Printf("не удалось преобразовать переменную PASSWORD_LENGTH_MIN в целое число: %v\n", err)
	}
	moodDescLenMax, err := strconv.Atoi(os.Getenv("MOOD_DESC_LENGTH_MAX"))
	if err != nil {
		fmt.Printf("не удалось преобразовать переменную MOOD_DESC_LENGTH_MAX в целое число: %v\n", err)
	}
	moodEmotesLenMax, err := strconv.Atoi(os.Getenv("MOOD_EMOTES_LENGTH_MAX"))
	if err != nil {
		fmt.Printf("не удалось преобразовать переменную MOOD_EMOTES_LENGTH_MAX в целое число: %v\n", err)
	}

	systemPrompt := `
Ты — заботливый помощник по ментальному здоровью. Тебе даётся история изменений настроения человека в виде массива объектов moods, а также previous_advice — прошлый совет, если он был. Каждый элемент moods содержит:

* score — уровень настроения (от 1 до 5),
* emotions — эмоции,
* description — комментарий от пользователя (может быть пустым),
* date — дата записи.

Твоя задача — на основе последних 3–5 записей выдать короткий, но полезный совет, который поможет человеку улучшить или сохранить хорошее настроение. Если previous_advice есть, избегай его повторения. Не нужно пересказывать данные — только вывод и совет. Совет должен быть кратким, в пределах 2–3 предложений.

Пример структуры входных данных:

{
    "previous_advice": "",
    "moods": [
        {
            "score": 5,
            "emotions": "",
            "description": "",
            "date": "09.05.2025"
        },
        {
            "score": 5,
            "emotions": "",
            "description": "",
            "date": "10.05.2025"
        }
    ]
}

Твой ответ:
Краткий совет по улучшению настроения.
`

	return &Config{
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),

		JWT_COOKIE_NAME: os.Getenv("JWT_COOKIE_NAME"),
		JWT_SECRET:      os.Getenv("JWT_SECRET"),

		GOOGLE_CLIENT_ID:       os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		GOOGLE_CLIENT_CALLBACK: os.Getenv("GOOGLE_CLIENT_CALLBACK"),

		GITHUB_CLIENT_ID:       os.Getenv("GITHUB_CLIENT_ID"),
		GITHUB_CLIENT_SECRET:   os.Getenv("GITHUB_CLIENT_SECRET"),
		GITHUB_CLIENT_CALLBACK: os.Getenv("GITHUB_CLIENT_CALLBACK"),

		SYSTEM_PROMPT: systemPrompt,
		AI_API_KEY:    os.Getenv("AI_API_KEY"),

		PASSWORD_LENGTH_MIN:    passwordLenMin,
		MOOD_DESC_LENGTH_MAX:   moodDescLenMax,
		MOOD_EMOTES_LENGTH_MAX: moodEmotesLenMax,

		JWT_HTTP_ONLY: false,
		JWT_SECURE:    false,
	}
}

//  = 5432
//  = "access_token"
