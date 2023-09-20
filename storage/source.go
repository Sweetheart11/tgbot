package storage

import "github.com/jmoiron/sqlx"

type SourcePostgresStorage struct {
	db *sqlx.DB
}
