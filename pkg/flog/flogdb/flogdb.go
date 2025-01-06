package flogdb

import (
	"database/sql"
	"io"
	"os"
	"path"
)

func InitSchema(db *sql.DB) error {
	filename := path.Join("schema", "schema.sql")
	dbFile, err := os.Open(filename)

	if err != nil {
		return err
	}

	dbSchema, err := io.ReadAll(dbFile)

	if err != nil {
		return err
	}

	_, err = db.Exec(string(dbSchema))

	if err != nil {
		return err
	}

	return nil
}
