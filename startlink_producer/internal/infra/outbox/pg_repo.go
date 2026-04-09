package outbox

import (
	"context"
	"errors"
	"time"

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

// Save пишет событие в той же транзакции, что и основной запрос (берёт tx из ctx).
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

// FetchUnprocessed блокирует строки (SKIP LOCKED) — безопасно при нескольких репликах
func (r *PgOutboxRepo) FetchUnprocessed(ctx context.Context, limit int) ([]outbox.Event, error) {
	query, args, err := r.sq.
		Select("id", "event_type", "payload", "created_at").
		From("outbox").
		Where(sq.Eq{"processed_at": nil}).
		OrderBy("created_at").
		Limit(uint64(limit)).
		Suffix("FOR UPDATE SKIP LOCKED").
		ToSql()
	if err != nil {
		return nil, err
	}
	runner := db.RunnerFromCtx(ctx, r.dbConn)
	rows, err := runner.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []outbox.Event
	for rows.Next() {
		var e outbox.Event
		if err := rows.Scan(&e.ID, &e.EventType, &e.Payload, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *PgOutboxRepo) MarkProcessed(ctx context.Context, id string) error {
	query, args, err := r.sq.
		Update("outbox").
		Set("processed_at", time.Now()).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	runner := db.RunnerFromCtx(ctx, r.dbConn)
	result, err := runner.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return errors.New("outbox event not found")
	}
	return nil
}
