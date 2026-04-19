package rest

import (
	"encoding/json"
	"net/http"

	"starlink_producer/domain/users"
)

type UserHandler struct {
	usecase users.UserUsecase
}

func NewUserHandler(usecase users.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req users.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.usecase.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
