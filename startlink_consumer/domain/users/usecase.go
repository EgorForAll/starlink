package users

import (
	"context"
	"encoding/json"
	"fmt"
)

type UserUsecase interface {
	Handle(ctx context.Context, payload []byte) error
}

// Transactor — интерфейс управления транзакцией
type Transactor interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Usecase struct {
	txManager Transactor
	service   *UserService
	repo      UserRepo
}

func NewUsecase(txManager Transactor, service *UserService, repo UserRepo) *Usecase {
	return &Usecase{
		txManager: txManager,
		service:   service,
		repo:      repo,
	}
}

// Handle декодирует payload из Kafka, обрабатывает данные и атомарно сохраняет в БД.
// Kafka offset коммитится снаружи только после успешного возврата этого метода.
func (uc *Usecase) Handle(ctx context.Context, payload []byte) error {
	var user User
	if err := json.Unmarshal(payload, &user); err != nil {
		return fmt.Errorf("unmarshal user event: %w", err)
	}

	received := uc.service.Process(user)

	return uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		if err := uc.repo.Save(ctx, &received); err != nil {
			return fmt.Errorf("save received user: %w", err)
		}
		return nil
	})
}
