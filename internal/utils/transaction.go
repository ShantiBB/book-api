package utils

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CloseTransaction(ctx context.Context, tx pgx.Tx, err error) {
	if p := recover(); p != nil {
		_ = tx.Rollback(ctx)
		panic(p)
	} else if err != nil {
		_ = tx.Rollback(ctx)
	} else {
		_ = tx.Commit(ctx)
	}
}
