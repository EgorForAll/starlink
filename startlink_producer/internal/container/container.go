package container

import (
	"os"

	"starlink_producer/domain/users"
	"starlink_producer/internal/config"
	"starlink_producer/internal/infra/db"
	infraoutbox "starlink_producer/internal/infra/outbox"

	"github.com/rs/zerolog"
)

type DiContainer struct {
	DbConn    db.DbConn
	Logger    zerolog.Logger
	HttpPort  string
	TxManager *db.TxManager

	UserUsecase users.UserUsecase
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

	c.TxManager = db.NewTxManager(c.DbConn)
	pgOutboxRepo := infraoutbox.NewPgOutboxRepo(c.DbConn)
	pgUserRepo := users.NewUserRepo(c.DbConn)

	userService := users.NewUserService(pgUserRepo)
	c.UserUsecase = users.NewUsecase(c.TxManager, userService, pgOutboxRepo)
}
