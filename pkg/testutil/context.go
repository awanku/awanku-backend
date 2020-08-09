package testutil

import (
	"context"
	"os"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/go-pg/pg/v9"
)

type key string

var ormKey key = "testutil_orm"

// Context creates context for test
func Context() (context.Context, func()) {
	db, closeDB := DBCluster()
	ctx := context.WithValue(context.Background(), appctx.KeyDatabase, db)

	opts, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("invalid DATABASE_URL")
	}
	conn := pg.Connect(opts)
	ctx = context.WithValue(ctx, ormKey, conn)

	cleanFn := func() {
		closeDB()
		conn.Close()
	}
	return ctx, cleanFn
}

func orm(ctx context.Context) *pg.DB {
	raw := ctx.Value(ormKey)
	if val, ok := raw.(*pg.DB); ok {
		return val
	}
	return nil
}
