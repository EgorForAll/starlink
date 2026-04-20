package app

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"starlink_consumer/internal/config"
	"starlink_consumer/internal/container"
	"starlink_consumer/internal/infra/alerting"

	"github.com/segmentio/kafka-go"
)

var retryDelays = []time.Duration{5 * time.Second, 30 * time.Second, 5 * time.Minute}

func InitApp(di *container.DiContainer, cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	retryTopics := []string{
		cfg.KafkaTopic + ".retry.1",
		cfg.KafkaTopic + ".retry.2",
		cfg.KafkaTopic + ".retry.3",
	}
	dlqTopic := cfg.KafkaTopic + ".dlq"

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBrokers...),
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	router := NewRouter(writer, alerting.NewStubAlerter(di.Logger), retryTopics, dlqTopic, di.Logger)

	errCh := make(chan error, 1+len(retryTopics))

	launch := func(run func(context.Context) error) {
		go func() {
			if err := run(ctx); err != nil {
				errCh <- err
			}
		}()
	}

	di.Logger.Info().Msg("starting kafka consumer")
	launch(NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID, di.UserUsecase, router, di.Logger).Run)

	for i, topic := range retryTopics {
		rc := NewRetryConsumer(cfg.KafkaBrokers, topic, cfg.KafkaGroupID, retryDelays[i], di.UserUsecase, router, di.Logger)
		di.Logger.Info().Str("topic", topic).Dur("delay", retryDelays[i]).Msg("starting retry consumer")
		launch(rc.Run)
	}

	select {
	case <-ctx.Done():
		di.Logger.Info().Msg("shutting down consumers")
	case err := <-errCh:
		di.Logger.Error().Err(err).Msg("consumer error")
		return err
	}

	return nil
}
