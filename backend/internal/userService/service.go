package userService

import (
	"errors"
	"sentimenta/internal/hash"
	"sentimenta/internal/utils"
)

type UserService interface {
	CreateUser(username string, email string, password string) (User, error)
	GetUser(id string) (User, error)
	UpdateUser(u UserUpdate) (User, error)
	DeleteUser(id string) error
	Authenticate(email string, password string) (User, error)
}

type userService struct {
	repo UserRepository
}

func (s *userService) CreateUser(username string, email string, password string) (User, error) {
	if !utils.IsValidEmail(email) {
		return User{}, errors.New("email не прошел валидацию")
	}

	passwordHash := hash.HashPassword(password)

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

func (s *userService) UpdateUser(u UserUpdate) (User, error) {
	targetUser, err := s.repo.GetUser(u.Uid)
	if err != nil {
		return User{}, err
	}

	if u.Password != nil {
		passwordHash := hash.HashPassword(*u.Password)
		targetUser.PasswordHash = passwordHash
	}
	if u.Username != nil {
		targetUser.Username = *u.Username
	}
	if u.Email != nil {
		targetUser.Email = *u.Email
	}

	if err := s.repo.UpdateUser(targetUser); err != nil {
		return User{}, err
	}
	return targetUser, nil
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) Authenticate(email string, password string) (User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return User{}, err
	}
	if !hash.VerifyPassword(password, user.PasswordHash) {
		return User{}, errors.New("Неверный пароль")
	}
	return user, nil
}

func NewService(r UserRepository) UserService {
	return &userService{repo: r}
}
