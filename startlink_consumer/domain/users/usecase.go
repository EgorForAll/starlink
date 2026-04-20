package users

import (
	"context"
	"encoding/json"
	"fmt"
)

type UserUsecase interface {
	Handle(ctx context.Context, payload []byte) error
}

type Usecase struct {
	service *UserService
	repo    UserRepo
}

func NewUsecase(service *UserService, repo UserRepo) *Usecase {
	return &Usecase{
		service: service,
		repo:    repo,
	}
}

// Handle декодирует payload из Kafka, обрабатывает данные и сохраняет в БД.
// Ошибки разделены: невалидный payload — NonRetryableError, сбой БД — RetryableError.
func (uc *Usecase) Handle(ctx context.Context, payload []byte) error {
	var user User
	if err := json.Unmarshal(payload, &user); err != nil {
		return &NonRetryableError{Err: fmt.Errorf("unmarshal user event: %w", err)}
	}

	received := uc.service.Process(user)

	if err := uc.repo.Save(ctx, &received); err != nil {
		return &RetryableError{Err: fmt.Errorf("save received user: %w", err)}
	}
	return nil
}
