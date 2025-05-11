package repository

import (
	m "sentimenta/internal/models"
	"time"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type AdviceRepository interface {
	GetAdvices(userID string) ([]m.Advice, error)
	GetAdvice(userID string, date time.Time) (m.Advice, error)
	CreateAdvice(a *m.Advice) error
	GetLastAdvice(userID string) (m.Advice, error)
}

type MoodRepository interface {
	GetMoods(userID string) ([]m.Mood, error)
	GetLastMoods(userID string, limit int) ([]m.Mood, error)
	CreateMood(m *m.Mood) error
	UpdateMood(m *m.Mood) error
	DeleteMood(id string) error
}

type UserRepository interface {
	CreateUser(user *m.User) error
	GetUser(id string) (m.User, error)
	GetAllUsers() ([]m.User, error)
	GetUserByEmail(email string) (*m.User, error)
	UpdateUser(userID int, updates any) error
	DeleteUser(id string) error
}
