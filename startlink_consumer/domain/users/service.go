package users

import "strings"

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// Process имитирует обработку данных:
// переводит имя и фамилию в верхний регистр, нормализует email
func (s *UserService) Process(user User) ReceivedUser {
	return ReceivedUser{
		SourceID:  user.ID,
		FirstName: strings.ToUpper(user.FirstName),
		LastName:  strings.ToUpper(user.LastName),
		Email:     strings.ToLower(user.Email),
	}
}
