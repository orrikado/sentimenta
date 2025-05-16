package service

import (
	"errors"
	"fmt"
	errs "sentimenta/internal/errors"
	"sentimenta/internal/hash"
	m "sentimenta/internal/models"
	repo "sentimenta/internal/repository"
	"sentimenta/internal/utils"

	"gorm.io/gorm"
)

type userService struct {
	repo repo.UserRepository
}

func (s *userService) CreateUser(username string, email string, password *string, timezone string) (m.User, error) {
	if !utils.IsValidEmail(email) {
		return m.User{}, errs.ErrEmailValidation
	}

	existingUser, err := s.repo.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		return *existingUser, errs.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return m.User{}, err
	}

	var passwordHashPtr *string
	if password != nil {
		hashed := hash.HashPassword(*password)
		passwordHashPtr = &hashed
	}

	newUser := m.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHashPtr,
		Timezone:     timezone,
	}
	if err := s.repo.CreateUser(&newUser); err != nil {
		return m.User{}, err
	}

	return newUser, nil
}

func (s *userService) GetUser(id string) (m.User, error) {
	return s.repo.GetUser(id)
}

func (s *userService) UpdateUser(userID string, r m.UserUpdateReq) (m.User, error) {
	targetUser, err := s.repo.GetUser(userID)
	if err != nil {
		return m.User{}, err
	}

	updates := map[string]any{}

	if r.Username != nil {
		updates["username"] = *r.Username
	}
	if r.Email != nil {
		updates["email"] = *r.Email
	}
	if r.UseAI != nil {
		updates["use_ai"] = *r.UseAI
	}
	if r.Timezone != nil {
		updates["timezone"] = *r.Timezone
	}

	if len(updates) == 0 {
		return m.User{}, nil
	}

	if err := s.repo.UpdateUser(targetUser.Uid, updates); err != nil {
		return m.User{}, err
	}
	return targetUser, nil
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) Authenticate(email, password string) (m.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return m.User{}, err
	}
	if !hash.VerifyPassword(password, *user.PasswordHash) {
		return m.User{}, errs.ErrWrongPassword
	}
	fmt.Printf("Authenticate m.User: %v", user)

	return *user, nil
}

func (s *userService) ChangePassword(userID, password, newPassword string) error {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return err
	}

	newPasswordHash := hash.HashPassword(newPassword)
	if !hash.VerifyPassword(newPasswordHash, *user.PasswordHash) {
		return errs.ErrWrongPassword
	}

	passwordHash := hash.HashPassword(newPassword)
	user.PasswordHash = &passwordHash

	if err := s.repo.UpdateUser(user.Uid, user); err != nil {
		return err
	}
	return nil
}

func (s *userService) GetUserByEmail(email string) (m.User, error) {
	result, err := s.repo.GetUserByEmail(email)
	return *result, err
}

func NewUserService(r repo.UserRepository) UserService {
	return &userService{repo: r}
}
