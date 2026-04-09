package rest

import (
	"net/http"
	"github.com/rs/zerolog"
)

func NewRouter(logger zerolog.Logger, userHandler *UserHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)

	return Logger(logger)(mux)
}