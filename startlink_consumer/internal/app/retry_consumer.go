package app

import (
	"context"
	"time"

	"starlink_consumer/domain/users"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type RetryConsumer struct {
	reader  *kafka.Reader
	usecase users.UserUsecase
	router  *Router
	delay   time.Duration
	logger  zerolog.Logger
}

func NewRetryConsumer(brokers []string, topic, groupID string, delay time.Duration, usecase users.UserUsecase, router *Router, logger zerolog.Logger) *RetryConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &RetryConsumer{
		reader:  reader,
		usecase: usecase,
		router:  router,
		delay:   delay,
		logger:  logger,
	}
}

// Run читает retry-сообщения, выдерживает delay перед обработкой, затем роутит при ошибке
func (c *RetryConsumer) Run(ctx context.Context) error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil // graceful shutdown
			}
			c.logger.Error().Err(err).Msg("kafka: fetch message error")
			return err
		}

		c.logger.Info().
			Str("topic", msg.Topic).
			Int64("offset", msg.Offset).
			Dur("delay", c.delay).
			Msg("kafka: retry message received, waiting before processing")

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(c.delay):
		}

		if err := c.usecase.Handle(ctx, msg.Value); err != nil {
			c.router.Route(ctx, msg, err)
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error().Err(err).Msg("kafka: commit message error")
			return err
		}

		c.logger.Info().Int64("offset", msg.Offset).Msg("kafka: retry message committed")
	}
}
