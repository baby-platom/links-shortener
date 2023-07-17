package database

import (
	"context"
	"database/sql"
)

// CreateShortIDsTable creates short_ids table
func (db *DB) CreateShortIDsTable(ctx context.Context) error {
	_, err := db.connection.ExecContext(
		ctx,
		`CREATE TABLE short_ids("id" TEXT, "url" TEXT);`,
	)
	return err
}

func (db *DB) CheckIfShortIDsTableExists(ctx context.Context) (bool, error) {
	row := db.connection.QueryRowContext(
		ctx,
		`SELECT EXISTS (
            SELECT *
            FROM information_schema.tables
            WHERE
              table_name = 'short_ids'
        );`,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

func (db *DB) WriteShortenedURL(ctx context.Context, id string, url string) error {
	_, err := db.connection.ExecContext(
		ctx,
		"INSERT INTO short_ids VALUES($1,$2);",
		id,
		url,
	)
	return err
}

func (db *DB) GetShortenedURL(ctx context.Context, id string) (string, error) {
	row := db.connection.QueryRowContext(
		ctx,
		"SELECT url FROM short_ids WHERE id=$1;",
		id,
	)
	var url string
	err := row.Scan(&url)
	if err == sql.ErrNoRows {
		err = nil
	}
	return url, err
}
