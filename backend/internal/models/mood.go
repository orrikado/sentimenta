package models

import (
	"time"
)

type Mood struct {
	Uid         int       `json:"uid" gorm:"primaryKey;autoIncrement;unique"`
	Score       int16     `json:"score" gorm:"type:SMALLINT"`
	Emotions    string    `json:"emotions"`
	Description string    `json:"description"`
	UserId      int       `json:"user_id"`
	Date        time.Time `json:"date" gorm:"type:date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MoodAdd struct {
	Score       int16     `json:"score"`
	Emotions    string    `json:"emotions"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
}

type MoodUpdate struct {
	Uid         int     `json:"uid"`
	Score       *int16  `json:"score,omitempty"`
	Emotions    *string `json:"emotions,omitempty"`
	Description *string `json:"description,omitempty"`
}

type MoodDTO struct {
	MoodAdd
}
