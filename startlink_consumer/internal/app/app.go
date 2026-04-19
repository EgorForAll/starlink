package app

import (
	"context"
	"os/signal"
	"syscall"

	"starlink_consumer/internal/config"
	"starlink_consumer/internal/container"
)

func InitApp(di *container.DiContainer, cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	consumer := NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID, di.UserUsecase, di.Logger)

	errCh := make(chan error, 1)

	go func() {
		di.Logger.Info().Msg("starting kafka consumer")
		if err := consumer.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		di.Logger.Info().Msg("shutting down consumer")
	case err := <-errCh:
		di.Logger.Error().Err(err).Msg("consumer error")
		return err
	}

	return nil
}
