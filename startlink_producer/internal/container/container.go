package container

import (
	"os"
	"time"

	"starlink_producer/domain/users"
	"starlink_producer/internal/config"
	"starlink_producer/internal/infra/db"
	"starlink_producer/internal/infra/kafka"
	infraoutbox "starlink_producer/internal/infra/outbox"

	"github.com/rs/zerolog"
)

type DiContainer struct {
	DbConn    db.DbConn
	Logger    zerolog.Logger
	HttpPort  string
	TxManager *db.TxManager

	// usecase
	UserUsecase users.UserUsecase

	// relay
	OutboxRelay *infraoutbox.Relay
}

func NewDiContainer(dbConn db.DbConn, cfg *config.Config) *DiContainer {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("app", cfg.AppName).Logger()
	return &DiContainer{
		DbConn:   dbConn,
		Logger:   logger,
		HttpPort: cfg.HttpPort,
	}
}

func (c *DiContainer) InitDependencies(cfg *config.Config) {
	if c.DbConn == nil {
		panic("db is not initialized")
	}

	// infra
	c.TxManager = db.NewTxManager(c.DbConn)
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers)
	pgOutboxRepo := infraoutbox.NewPgOutboxRepo(c.DbConn)
	pgUserRepo := users.NewUserRepo(c.DbConn)

	// domain
	userService := users.NewUserService(pgUserRepo)
	c.UserUsecase = users.NewUsecase(c.TxManager, userService, pgOutboxRepo)

	// relay
	c.OutboxRelay = infraoutbox.NewRelay(
		c.TxManager,
		pgOutboxRepo,
		kafkaProducer,
		5*time.Second,
		c.Logger,
	)
}
