package failure

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type PgErrorCode = string

const (
	ErrPgUniqueViolation PgErrorCode = "23505"
)

func PgError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}
