package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(parentCtx context.Context, user *User) error
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(parentCtx context.Context, user User) (*User, error) {
	existedUser, err := s.repo.FindByEmail(parentCtx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if existedUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	if err = s.repo.Create(parentCtx, &user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}
