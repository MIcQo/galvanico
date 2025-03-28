package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(_ context.Context, _ *bun.DB) error {
		return nil
	}, func(_ context.Context, _ *bun.DB) error {
		return nil
	})
}
