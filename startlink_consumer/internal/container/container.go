package container

import (
	"os"

	"starlink_consumer/domain/users"
	"starlink_consumer/internal/config"
	"starlink_consumer/internal/infra/db"

	"github.com/rs/zerolog"
)

type DiContainer struct {
	DbConn      db.DbConn
	Logger      zerolog.Logger
	TxManager   *db.TxManager
	UserUsecase users.UserUsecase
}

func NewDiContainer(dbConn db.DbConn, cfg *config.Config) *DiContainer {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("app", cfg.AppName).Logger()
	return &DiContainer{
		DbConn: dbConn,
		Logger: logger,
	}
}

func (c *DiContainer) InitDependencies() {
	if c.DbConn == nil {
		panic("db is not initialized")
	}

	c.TxManager = db.NewTxManager(c.DbConn)
	pgUserRepo := users.NewUserRepo(c.DbConn)
	userService := users.NewUserService()
	c.UserUsecase = users.NewUsecase(c.TxManager, userService, pgUserRepo)
}
