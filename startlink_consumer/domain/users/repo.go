package users

import (
	"context"
	"time"

	"starlink_consumer/internal/infra/db"

	sq "github.com/Masterminds/squirrel"
)

type UserRepo interface {
	Save(ctx context.Context, user *ReceivedUser) error
}

type PgUserRepo struct {
	dbConn db.DbConn
	sq     sq.StatementBuilderType
}

func NewUserRepo(dbConn db.DbConn) *PgUserRepo {
	return &PgUserRepo{
		dbConn: dbConn,
		sq:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PgUserRepo) runner(ctx context.Context) db.TxRunner {
	return db.RunnerFromCtx(ctx, r.dbConn)
}

func (r *PgUserRepo) Save(ctx context.Context, user *ReceivedUser) error {
	query, args, err := r.sq.
		Insert("received_users").
		Columns("user_id", "first_name", "last_name", "email", "processed_at").
		Values(user.UserID, user.FirstName, user.LastName, user.Email, time.Now()).
		Suffix("RETURNING id, received_at").
		ToSql()
	if err != nil {
		return err
	}

	return r.runner(ctx).QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.ReceivedAt)
}
