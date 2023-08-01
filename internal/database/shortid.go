package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/baby-platom/links-shortener/internal/models"
)

var insertTemplate = "INSERT INTO short_ids (id, url, user_id) VALUES($1,$2,$3);"

// CreateShortIDsTable creates short_ids table
func (db *DB) CreateShortIDsTable(ctx context.Context) error {
	tx, err := db.connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE short_ids("id" TEXT PRIMARY KEY, "url" TEXT, "user_id" INT);`,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `CREATE UNIQUE INDEX original_url ON short_ids (url)`)
	if err != nil {
		return err
	}

	return tx.Commit()
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

func (db *DB) WriteShortenedURL(ctx context.Context, id string, url string, userID int) error {
	_, err := db.connection.ExecContext(
		ctx,
		insertTemplate,
		id,
		url,
		userID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
	}

	return err
}

func (db *DB) GetInitialURLByID(ctx context.Context, id string) (string, error) {
	row := db.connection.QueryRowContext(
		ctx,
		"SELECT url FROM short_ids WHERE id=$1",
		id,
	)
	var url string
	err := row.Scan(&url)
	if err == sql.ErrNoRows {
		err = nil
	}
	return url, err
}

func (db *DB) GetIDByInitialURL(ctx context.Context, url string) (string, error) {
	row := db.connection.QueryRowContext(
		ctx,
		"SELECT id FROM short_ids WHERE url=$1",
		url,
	)
	var id string
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return id, err
}

func (db *DB) WriteBatchOfShortenedURL(ctx context.Context, shortenedUrlsByIds []models.BatchPortionShortenResponse, userID int) error {
	tx, err := db.connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(
		ctx,
		insertTemplate,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, portion := range shortenedUrlsByIds {
		_, err := stmt.ExecContext(ctx, portion.ID, portion.OriginalURL, userID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (db *DB) GetUserShortenURLsList(ctx context.Context, baseAddress string, userID int) ([]models.UserShortenURLsListResponse, error) {
	result := make([]models.UserShortenURLsListResponse, 0)
	rows, err := db.connection.QueryContext(
		ctx,
		"SELECT id, url FROM short_ids WHERE user_id=$1;",
		userID,
	)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var piece models.UserShortenURLsListResponse
		err = rows.Scan(&piece.ShortURL, &piece.OriginalURL)
		if err != nil {
			return result, err
		}

		piece.ShortURL = fmt.Sprintf("%s/%s", baseAddress, piece.ShortURL)
		result = append(result, piece)
	}

	err = rows.Err()
	if err != nil {
		return result, err
	}
	return result, nil
}
