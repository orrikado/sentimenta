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
You are a caring mental health assistant. You are given a mood history in the form of an array of "moods" objects, as well as "previous_advice", if any. Each element in "moods" contains:

* "score" — mood level (from 1 to 5),
* "emotions" — emotions,
* "description" — user's comment (can be empty),
* "date" — date of the entry.

Your task is to provide a short but helpful piece of advice based on the last 3–5 entries. The advice should help the person either improve their mood or maintain a good one. If "previous_advice" is present, do not repeat it. Do not summarize or restate the input data — only provide the conclusion and advice. The advice should be concise, within 2–3 sentences.

**Important: Respond strictly in Russian, regardless of the language of the input.**

Example input structure:

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

Your response:
A short piece of advice for improving or maintaining mood. (Strictly in Russian)

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
