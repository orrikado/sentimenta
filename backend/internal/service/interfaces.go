package service

import (
	m "sentimenta/internal/models"
	"time"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type UserService interface {
	CreateUser(username, email string, password *string, timezone string) (m.User, error)
	GetUser(id string) (m.User, error)
	UpdateUser(userID string, u m.UserUpdateReq) (m.User, error)
	DeleteUser(id string) error
	ChangePassword(userID, password, newPassword string) error
	Authenticate(email, password string) (m.User, error)
	GetUserByEmail(email string) (m.User, error)
}

type MoodService interface {
	GetMoods(userID string) ([]m.Mood, error)
	CreateMood(userID string, score int16, emotions, description string, date time.Time) (m.Mood, error)
	UpdateMood(userID string, m *m.Mood) error
	DeleteMood(id string) error
}

type AdviceService interface {
	GetAdvice(userID string, date time.Time) (m.Advice, error)
	GetAdvices(userID string) ([]m.Advice, error)
	CreateAdvice(userID string, text string, date time.Time) (m.Advice, error)
	GetLastAdvice(userID string) (m.Advice, error)
	GenerateAdvice(userID int, date time.Time) (m.Advice, error)
}
