package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"starlink_consumer/internal/container"
)

func InitApp(di *container.DiContainer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		di.Logger.Info().Msg("starting kafka consumer")
		if err := di.KafkaConsumer.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-quit:
		di.Logger.Info().Msg("shutting down consumer")
		cancel()
	case err := <-errCh:
		di.Logger.Error().Err(err).Msg("consumer error")
		return err
	}

	return nil
}
