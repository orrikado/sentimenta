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
	AI_ENABLED    bool

	PASSWORD_LENGTH_MIN    int
	MOOD_DESC_LENGTH_MAX   int
	MOOD_EMOTES_LENGTH_MAX int

	JWT_HTTP_ONLY bool
	JWT_SECURE    bool

	REGISTRATION_ENABLED bool
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

You are a caring mental health assistant. You receive an "AdviceRequest" object containing:

* "previous_advice": a previous piece of advice, if any.
* "last_mood": the most recent mood entry.
* "moods": an array of previous mood entries (excluding "last_mood").

Each mood entry (both "last_mood" and items in "moods") has the following structure:

* "score" — mood level (from 1 to 5),
* "emotions" — emotions experienced,
* "description" — user's comment (can be empty),
* "date" — date of the entry.

Your task is to generate a short but helpful piece of advice that is **primarily based on the "last_mood" entry**, while **also lightly considering the general trend in "moods"**. If "previous_advice" is present, avoid repeating it. Do **not** summarize or restate the input data — only provide a direct and meaningful conclusion.

The advice should be concise (2–3 sentences) and supportive — aimed at improving or maintaining the person's emotional well-being.

**Important: Respond strictly in Russian, regardless of the language of the input.**

Example input structure:


{
  "previous_advice": "Старайся больше гулять на свежем воздухе.",
  "last_mood": {
    "score": 2,
    "emotions": "усталость, тревога",
    "description": "Много дел, ничего не успеваю.",
    "date": "21.05.2025"
  },
  "moods": [
    {
      "score": 3,
      "emotions": "раздражение",
      "description": "Сложный день.",
      "date": "20.05.2025"
    },
    {
      "score": 4,
      "emotions": "спокойствие",
      "description": "",
      "date": "19.05.2025"
    }
  ]
}

Your response:
A brief and helpful piece of advice in Russian, based mainly on "last_mood".

`

	return &Config{
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),

		JWT_COOKIE_NAME: os.Getenv("JWT_COOKIE_NAME"),
		JWT_SECRET:      os.Getenv("JWT_SECRET"),

		GOOGLE_CLIENT_ID:       os.Getenv("PUBLIC_GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		GOOGLE_CLIENT_CALLBACK: os.Getenv("GOOGLE_CLIENT_CALLBACK"),

		GITHUB_CLIENT_ID:       os.Getenv("PUBLIC_GITHUB_CLIENT_ID"),
		GITHUB_CLIENT_SECRET:   os.Getenv("GITHUB_CLIENT_SECRET"),
		GITHUB_CLIENT_CALLBACK: os.Getenv("GITHUB_CLIENT_CALLBACK"),

		SYSTEM_PROMPT: systemPrompt,
		AI_API_KEY:    os.Getenv("AI_API_KEY"),
		AI_ENABLED:    os.Getenv("AI_ENABLED") == "true",

		PASSWORD_LENGTH_MIN:    passwordLenMin,
		MOOD_DESC_LENGTH_MAX:   moodDescLenMax,
		MOOD_EMOTES_LENGTH_MAX: moodEmotesLenMax,

		JWT_HTTP_ONLY: false,
		JWT_SECURE:    false,

		REGISTRATION_ENABLED: true,
	}
}

//  = 5432
//  = "access_token"
