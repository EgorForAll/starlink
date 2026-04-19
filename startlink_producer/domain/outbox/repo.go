package outbox

import "context"

type Repo interface {
	Save(ctx context.Context, event Event) error
}
