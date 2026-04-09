package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"starlink_producer/internal/container"
	"starlink_producer/internal/infra/transport/rest"
)

func InitApp(di *container.DiContainer) error {
	userHandler := rest.NewUserHandler(di.UserUsecase)
	router := rest.NewRouter(di.Logger, userHandler)

	srv := &http.Server{
		Addr:              ":" + di.HttpPort,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Outbox relay — фоновый воркер, отправляет события из outbox в Kafka
	go di.OutboxRelay.Run(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		di.Logger.Info().Str("addr", srv.Addr).Msg("starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			di.Logger.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	<-quit
	di.Logger.Info().Msg("shutting down server")

	cancel() // останавливаем relay

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	return srv.Shutdown(shutdownCtx)
}
