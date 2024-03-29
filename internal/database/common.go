package database

import (
	"context"
	"database/sql"
	"errors"
)

// ErrConflict - integrity constraint violation
var ErrConflict = errors.New("data conflict")

// DB represents a db connection
type DB struct {
	connection *sql.DB
	name       string
}

// Connection represents an opened DB
var Connection DB

// HealthCheck pings the db connection
func (db *DB) HealthCheck(ctx context.Context) bool {
	if (*db == DB{}) {
		return false
	}

	if err := db.connection.PingContext(ctx); err != nil {
		return false
	}
	return true
}

// Close closes db connection
func (db *DB) Close() error {
	return db.connection.Close()
}
