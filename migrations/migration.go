package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

//go:embed *.sql
var migs embed.FS

func init() {
	if err := Migrations.Discover(migs); err != nil {
		panic(err)
	}
}
