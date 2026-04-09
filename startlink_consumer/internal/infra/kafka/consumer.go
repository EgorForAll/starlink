package kafka

import (
	"context"

	"starlink_consumer/domain/users"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	usecase users.UserUsecase
	logger  zerolog.Logger
}

func NewConsumer(brokers []string, topic, groupID string, usecase users.UserUsecase, logger zerolog.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{
		reader:  reader,
		usecase: usecase,
		logger:  logger,
	}
}

// Run читает сообщения из Kafka в цикле.
// FetchMessage не коммитит offset — CommitMessages вызывается только после успешной записи в БД.
func (c *Consumer) Run(ctx context.Context) error {
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
			Int("partition", msg.Partition).
			Int64("offset", msg.Offset).
			Msg("kafka: received message")

		if err := c.usecase.Handle(ctx, msg.Value); err != nil {
			// не коммитим offset — при рестарте сообщение придёт повторно
			c.logger.Error().Err(err).Msg("kafka: handle message error, skipping commit")
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error().Err(err).Msg("kafka: commit message error")
			return err
		}

		c.logger.Info().Int64("offset", msg.Offset).Msg("kafka: message committed")
	}
}
