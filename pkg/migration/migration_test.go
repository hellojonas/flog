package migration

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func _TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db")

	if err != nil {
		t.Fatal(err)
	}

	if err := Migrate(db, "../../migrations"); err != nil {
		t.Fatal(err)
	}
}
