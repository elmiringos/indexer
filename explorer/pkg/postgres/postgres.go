package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/elmiringos/indexer/explorer/config"

	// Register the PostgreSQL driver
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Connection struct {
	db  *sql.DB
	log *zap.Logger
}

func NewPostgresConnection(cfg *config.Config, log *zap.Logger) *Connection {
	log.Debug("Connect to database", zap.String("url", cfg.PG.URL))
	db, err := sql.Open("postgres", cfg.PG.URL)

	if err != nil {
		log.Fatal("Error in opening db", zap.Error(err))
	}

	pc := &Connection{db: db, log: log}

	pc.Health()

	return pc
}

func (pc *Connection) Health() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := pc.db.PingContext(ctx)
	if err != nil {
		pc.log.Fatal("db down", zap.Error(err))
	}

	return true
}

func (pc *Connection) GetDb() *sql.DB {
	return pc.db
}
