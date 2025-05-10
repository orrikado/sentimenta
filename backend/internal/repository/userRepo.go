package repository

import (
	m "sentimenta/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *m.User) error
	GetUser(id string) (m.User, error)
	GetAllUsers() ([]m.User, error)
	GetUserByEmail(email string) (*m.User, error)
	UpdateUser(user m.User) error
	DeleteUser(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) GetUser(id string) (m.User, error) {
	var user m.User
	if err := r.db.First(&user, "uid = ?", id).Error; err != nil {
		return m.User{}, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(user *m.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user m.User) error {
	return r.db.Save(&user).Error
}

func (r *userRepository) DeleteUser(id string) error {
	return r.db.Delete(&m.User{}, "uid = ?", id).Error
}

func (r *userRepository) GetUserByEmail(email string) (*m.User, error) {
	var user m.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return &m.User{}, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]m.User, error) {
	var users []m.User
	if err := r.db.Find(&users).Error; err != nil {
		return []m.User{}, err
	}
	return users, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
