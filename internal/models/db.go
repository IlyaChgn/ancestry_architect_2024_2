package models

import "github.com/jackc/pgx/v4"

type EmptyRow struct{}

func (r EmptyRow) Scan(...interface{}) error {
	return pgx.ErrNoRows
}
