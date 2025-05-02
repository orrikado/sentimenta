package userservice

import "gorm.io/gorm"

type UserRepository interface {
	CreateUser(user User) error
	GetUser(id string) (User, error)
	UpdateUser(user User) error
	DeleteUser(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) GetUser(id string) (User, error) {
	var user User
	err := r.db.First(&user, "id = ?", id).Error
	return user, err
}

func (r *userRepository) CreateUser(user User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user User) error {
	return r.db.Save(&user).Error
}

func (r *userRepository) DeleteUser(id string) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}

func NewRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
