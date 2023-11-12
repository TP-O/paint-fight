package pg

import (
	pggenerated "client/infra/persistence/pg/generated"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*pggenerated.Queries
	db *pgxpool.Pool
}

func (s *Store) Close() {
	s.db.Close()
}

// TODO: use tx
func (s *Store) execTx(ctx context.Context, fn func(q *pggenerated.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := pggenerated.New(tx)
	if err = fn(q); err != nil {
		if err = tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}
