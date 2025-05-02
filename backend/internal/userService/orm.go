package userService

import (
	m "sentimenta/internal/moodService"
	"time"
)

type User struct {
	Uid          int       `json:"uid" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Moods        []m.Mood  `json:"moods"`
}

type UserUpdate struct {
	Uid      string  `json:"uid" gorm:"primaryKey;autoIncrement"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
