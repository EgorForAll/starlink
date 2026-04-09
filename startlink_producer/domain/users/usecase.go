package users

import (
	"context"
	"encoding/json"
	"fmt"

	"starlink_producer/domain/outbox"
)

type UserUsecase interface {
	Create(ctx context.Context, user User) (*User, error)
}

// Transactor — интерфейс управления транзакцией
type Transactor interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Usecase struct {
	txManager  Transactor
	service    *UserService
	outboxRepo outbox.Repo
}

func NewUsecase(txManager Transactor, service *UserService, outboxRepo outbox.Repo) *Usecase {
	return &Usecase{
		txManager:  txManager,
		service:    service,
		outboxRepo: outboxRepo,
	}
}

// Create создаёт пользователя и атомарно записывает событие в outbox
func (uc *Usecase) Create(ctx context.Context, user User) (*User, error) {
	var created *User

	err := uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		var err error
		created, err = uc.service.Create(ctx, user)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		payload, err := json.Marshal(created)
		if err != nil {
			return fmt.Errorf("marshal user event: %w", err)
		}

		return uc.outboxRepo.Save(ctx, outbox.Event{
			EventType: "user.created",
			Payload:   payload,
		})
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

