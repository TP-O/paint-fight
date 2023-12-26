package pg

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hub/config"
	pggenerated "hub/infra/persistence/pg/generated"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *Store

func New(ctx context.Context, cfg config.PostgreSQL) *Store {
	sync.OnceFunc(func() {
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
		poolCfg.ConnConfig.TLSConfig = createTLSConfig(cfg.RootCA, cfg.Host)

		pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
		if err != nil {
			log.Fatalf("Unable to create connection pool: %s", err.Error())
		}

		if err := pool.Ping(ctx); err != nil {
			log.Fatalf("Unable to connect to database: %s", err.Error())
		}

		db = &Store{
			Queries: pggenerated.New(pool),
			db:      pool,
		}
	})()

	return db
}

func createTLSConfig(rootCAPath, serverName string) *tls.Config {
	root, err := os.ReadFile(rootCAPath)
	if err != nil {
		log.Fatalf("Unable to read Root CA file: %s", err.Error())
	}

	block, _ := pem.Decode(root)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Unable to parse certificate: %s", err.Error())
	}

	c := x509.NewCertPool()
	c.AddCert(cert)

	return &tls.Config{
		RootCAs:    c,
		ServerName: serverName,
	}
}
