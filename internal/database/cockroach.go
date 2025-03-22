package database

import (
	"context"
	cfg "galvanico/internal/config"
	"github.com/jackc/pgx/v5"
	"log"
	"sync"
)

var once sync.Once
var conn *pgx.Conn

func Connection() *pgx.Conn {
	once.Do(func() {
		var config, configErr = cfg.Load()
		if configErr != nil {
			log.Fatal(configErr)
		}

		connConfig, err := pgx.ParseConfig(config.Database.URL)
		if err != nil {
			log.Fatal(err)
		}

		connConfig.RuntimeParams["application_name"] = config.AppName

		conn, err = pgx.ConnectConfig(context.Background(), connConfig)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connected to database")
	})

	return conn
}

func Close() error {
	if conn == nil {
		return nil
	}

	return conn.Close(context.Background())
}
