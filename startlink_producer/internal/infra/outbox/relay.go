package outbox

import (
	"context"
	"time"

	"starlink_producer/domain/outbox"
	"starlink_producer/internal/infra/db"

	"github.com/rs/zerolog"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte) error
}

type Relay struct {
	txManager *db.TxManager
	repo      outbox.Repo
	publisher Publisher
	interval  time.Duration
	logger    zerolog.Logger
}

func NewRelay(
	txManager *db.TxManager,
	repo outbox.Repo,
	publisher Publisher,
	interval time.Duration,
	logger zerolog.Logger,
) *Relay {
	return &Relay{
		txManager: txManager,
		repo:      repo,
		publisher: publisher,
		interval:  interval,
		logger:    logger,
	}
}


func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.process(ctx); err != nil {
				r.logger.Error().Err(err).Msg("outbox relay: process error")
			}
		}
	}
}

// process читает необработанные события и публикует их в Kafka
func (r *Relay) process(ctx context.Context) error {
	return r.txManager.RunInTx(ctx, func(ctx context.Context) error {
		events, err := r.repo.FetchUnprocessed(ctx, 100)
		if err != nil || len(events) == 0 {
			return err
		}

		for _, event := range events {
			if err := r.publisher.Publish(ctx, event.EventType, event.Payload); err != nil {
				return err // rollback — попробуем на следующем тике
			}
			if err := r.repo.MarkProcessed(ctx, event.ID); err != nil {
				return err
			}
		}
		return nil
	})
}
