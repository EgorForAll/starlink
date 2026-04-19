package outbox

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"starlink_producer/domain/outbox"
	"starlink_producer/internal/infra/db"
)

type PgOutboxRepo struct {
	dbConn db.DbConn
	sq     sq.StatementBuilderType
}

func NewPgOutboxRepo(dbConn db.DbConn) *PgOutboxRepo {
	return &PgOutboxRepo{
		dbConn: dbConn,
		sq:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Save пишет событие в той же транзакции, что и основной запрос
// Debezium отслеживает новые строки через PostgreSQL WAL и публикует их в Kafka
func (r *PgOutboxRepo) Save(ctx context.Context, event outbox.Event) error {
	query, args, err := r.sq.
		Insert("outbox").
		Columns("event_type", "payload").
		Values(event.EventType, event.Payload).
		ToSql()
	if err != nil {
		return err
	}
	runner := db.RunnerFromCtx(ctx, r.dbConn)
	_, err = runner.ExecContext(ctx, query, args...)
	return err
}
