package userservice

import (
	"errors"
	"sentimenta/internal/utils"
)

type UserService interface {
	CreateUser(username string, email string, password string) (User, error)
	GetUser(id string) (User, error)
	UpdateUser(id string, username string, email string, password string) (User, error)
	DeleteUser(id string) error
}

type userService struct {
	repo UserRepository
}

func (s *userService) CreateUser(username string, email string, password string) (User, error) {
	if !utils.IsValidEmail(email) {
		return User{}, errors.New("email не прошел валидацию")
	}

	passwordHash := "" // !!! Сделать хеширование

	newUser := User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}
	if err := s.repo.CreateUser(newUser); err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (s *userService) GetUser(id string) (User, error) {
	return s.repo.GetUser(id)
}

func (s *userService) UpdateUser(id string, username string, email string, password string) (User, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return User{}, err
	}

	passwordHash := "" // !!! Сделать хеширование

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		user.PasswordHash = passwordHash
	}
	if err := s.repo.UpdateUser(user); err != nil {
		return User{}, err
	}
	return user, nil

}

func (s *userService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func NewService(r UserRepository) UserService {
	return &userService{repo: r}
}
