package moodservice

import (
	"time"
)

type Mood struct {
	Uid         int       `json:"uid" gorm:"primaryKey;autoIncrement"`
	Score       int       `json:"score"`
	Emotions    string    `json:"emotions"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" gorm:"type:date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
