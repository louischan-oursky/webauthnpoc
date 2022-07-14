package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDatabase() (*sql.DB, error) {
	return sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
}

func WithTx(ctx context.Context, db *sql.DB, f func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = f(tx)
	if err != nil {
		return
	}

	return
}
