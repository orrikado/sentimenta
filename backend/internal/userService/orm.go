package userservice

import "time"

type User struct {
	Uid          int       `json:"uid" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
