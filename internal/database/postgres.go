package database

import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"  // initialization
)

// OpenPostgres opens connection to postgres db
func OpenPostgres(dataSourceName string) error {
	db, err := sql.Open("pgx", dataSourceName)
    if err != nil {
        return err
    }
    Connection = DB{
        connection: db,
        name: "postgres",
    }
	return nil
}
