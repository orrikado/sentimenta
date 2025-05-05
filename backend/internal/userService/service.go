package userService

import (
	"errors"
	"fmt"
	errs "sentimenta/internal/errors"
	"sentimenta/internal/hash"
	"sentimenta/internal/utils"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(username, email string, password *string) (User, error)
	GetUser(id string) (User, error)
	UpdateUser(userID string, u UserUpdate) (User, error)
	DeleteUser(id string) error
	ChangePassword(userID, password, newPassword string) error
	Authenticate(email, password string) (User, error)
	GetUserByEmail(email string) (User, error)
}

type userService struct {
	repo UserRepository
}

func (s *userService) CreateUser(username string, email string, password *string) (User, error) {
	if !utils.IsValidEmail(email) {
		return User{}, errs.ErrEmailValidation
	}

	existingUser, err := s.repo.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		return *existingUser, errs.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, err
	}

	var passwordHashPtr *string
	if password != nil {
		hashed := hash.HashPassword(*password)
		passwordHashPtr = &hashed
	}

	newUser := User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHashPtr,
	}
	if err := s.repo.CreateUser(&newUser); err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (s *userService) GetUser(id string) (User, error) {
	return s.repo.GetUser(id)
}

func (s *userService) UpdateUser(userID string, u UserUpdate) (User, error) {
	targetUser, err := s.repo.GetUser(userID)
	if err != nil {
		return User{}, err
	}

	if u.Password != nil {
		passwordHash := hash.HashPassword(*u.Password)
		targetUser.PasswordHash = &passwordHash
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

func (s *userService) Authenticate(email, password string) (User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return User{}, err
	}
	if !hash.VerifyPassword(password, *user.PasswordHash) {
		return User{}, errs.ErrWrongPassword
	}
	fmt.Printf("Authenticate User: %v", user)

	return *user, nil
}

func (s *userService) ChangePassword(userID, password, newPassword string) error {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return err
	}

	if *user.PasswordHash != hash.HashPassword(password) {
		return errs.ErrWrongPassword
	}

	passwordHash := hash.HashPassword(newPassword)
	user.PasswordHash = &passwordHash

	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

func (s *userService) GetUserByEmail(email string) (User, error) {
	result, err := s.repo.GetUserByEmail(email)
	return *result, err
}

func NewService(r UserRepository) UserService {
	return &userService{repo: r}
}
