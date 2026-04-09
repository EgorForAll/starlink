package outbox

import "time"

type Event struct {
	ID          string
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
	ProcessedAt *time.Time
}
