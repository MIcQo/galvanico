package database

import (
	"context"
	"database/sql"
	cfg "galvanico/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
	"runtime"
	"sync"
)

var once sync.Once
var conn *bun.DB

// Connection create or get database connection
func Connection() *bun.DB {
	once.Do(func() {
		conn = createConnection()
	})

	return conn
}

// Close closes all connections
func Close() error {
	if conn == nil {
		return nil
	}

	log.Info().Msg("closing database")

	return conn.Close()
}

func createConnection() *bun.DB {
	var config, configErr = cfg.Load()
	if configErr != nil {
		log.Fatal().Err(configErr)
	}

	connConfig, err := createPgxConfig(config)
	if err != nil {
		log.Fatal().Err(err)
	}

	var sqldb = stdlib.OpenDB(*connConfig)
	var db = bun.NewDB(
		sqldb,
		pgdialect.New(),
		bun.WithDiscardUnknownColumns(),
	)

	setupMaxConnections(sqldb)
	registerHooks(db, config)

	log.Info().Msg("connected to database")

	return db
}

func registerHooks(db *bun.DB, cfg *cfg.Config) {
	db.AddQueryHook(addQueryTelemetry(db))
	db.AddQueryHook(addQueryLogger(cfg))
}

func addQueryLogger(config *cfg.Config) *bundebug.QueryHook {
	return bundebug.NewQueryHook(
		// disable the hook
		bundebug.WithEnabled(config.LogLevel == "debug"),
	)
}

func addQueryTelemetry(db *bun.DB) *bunotel.QueryHook {
	var database string
	if _, err := db.NewSelect().NewRaw("SELECT CURRENT_DATABASE()").Exec(context.Background(), &database); err != nil {
		log.Fatal().Err(err)
	}

	return bunotel.NewQueryHook(bunotel.WithDBName(database))
}

func setupMaxConnections(sqldb *sql.DB) {
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)
}

func createPgxConfig(config *cfg.Config) (*pgx.ConnConfig, error) {
	connConfig, err := pgx.ParseConfig(config.Database.URL)
	if err != nil {
		log.Fatal().Err(err)
	}
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	connConfig.RuntimeParams["application_name"] = config.AppName
	return connConfig, err
}
