package models

import (
	"time"
)

type Advice struct {
	Uid    int       `json:"uid" gorm:"primaryKey;autoIncrement;unique"`
	UserID int       `json:"user_id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date" gorm:"type:date"`
}

type AdviceRequest struct {
	PreviousAdvice string    `json:"previous_advice"`
	LastMood       MoodAdd   `json:"last_mood"`
	Moods          []MoodAdd `json:"moods"`
}
