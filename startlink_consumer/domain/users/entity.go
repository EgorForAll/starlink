package users

import "time"

//TODO: сделать в internal
// User — сообщение, которое приходит из топика user.created (от producer)
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ReceivedUser — запись в consumer_db после обработки
type ReceivedUser struct {
	ID          int
	UserID      int
	FirstName   string
	LastName    string
	Email       string
	ReceivedAt  time.Time
	ProcessedAt time.Time
}
