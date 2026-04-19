package users

import (
	"context"
	"time"

	"starlink_producer/internal/infra/db"

	"github.com/Masterminds/squirrel"
)

const dbQueryTimeout = 2 * time.Second

type PgUserRepo struct {
	dbConn db.DbConn
}

func NewUserRepo(dbConn db.DbConn) *PgUserRepo {
	return &PgUserRepo{dbConn: dbConn}
}

// runner возвращает *sql.Tx из контекста (если есть) либо dbConn.
// Это позволяет репозиторию участвовать в транзакции, начатой на уровне usecase.
func (r *PgUserRepo) runner(ctx context.Context) db.TxRunner {
	return db.RunnerFromCtx(ctx, r.dbConn)
}

func (r *PgUserRepo) FindByEmail(parentCtx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(parentCtx, dbQueryTimeout)
	defer cancel()

	qb := squirrel.
		Select("id", "first_name", "last_name", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = r.runner(ctx).QueryRowContext(ctx, query, args...).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PgUserRepo) Create(parentCtx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(parentCtx, dbQueryTimeout)
	defer cancel()

	qb := squirrel.
		Insert("users").
		Columns("first_name", "last_name", "email").
		Values(user.FirstName, user.LastName, user.Email).
		Suffix(`RETURNING id, created_at`).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	return r.runner(ctx).QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
	)
}
