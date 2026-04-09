package outbox

import "context"

type Repo interface {
	Save(ctx context.Context, event Event) error
	FetchUnprocessed(ctx context.Context, limit int) ([]Event, error)
	MarkProcessed(ctx context.Context, id string) error
}
