package userservice

import "gorm.io/gorm"

type UserRepository interface {
	GetUser(id string) (User, error)
	CreateUser(user User) error
	UpdateUser() error
	DeleteUser() error
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(id string) (User, error) {
	var user User
	err := r.db.First(&user, "id = ?", id).Error
	return user, err
}
