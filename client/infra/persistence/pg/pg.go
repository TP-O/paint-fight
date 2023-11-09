package pg

import (
	"client/config"
	pggenerated "client/infra/persistence/pg/generated"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func New(ctx context.Context, cfg config.PostgreSQL) *Store {
	if db == nil {
		poolCfg, err := pgxpool.ParseConfig(fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v?pool_max_conns=%d&sslmode=require",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.PoolSize,
		))
		if err != nil {
			log.Fatalf("Unable to parse connection string: %s", err.Error())
		}

		db, err = pgxpool.NewWithConfig(ctx, poolCfg)
		if err != nil {
			log.Fatalf("Unable to create connection pool: %s", err.Error())
		}
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %s", err.Error())
	}

	return &Store{
		Queries: pggenerated.New(db),
		db:      db,
	}
}
