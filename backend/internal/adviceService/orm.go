package adviceService

import (
	"sentimenta/internal/moodService"
	"time"
)

type Advice struct {
	UserID int       `json:"user_id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date" gorm:"type:date"`
}

type AdviceRequest struct {
	PreviousAdvice string                `json:"previous_advice"`
	Moods          []moodService.MoodAdd `json:"moods"`
}
