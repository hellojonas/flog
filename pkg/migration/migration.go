package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Migrate(db *sql.DB, srcPath string) error {
	db.Exec("PRAGMA foreign_keys = ON;")
	db.Exec("CREATE TABLE IF NOT EXISTS __migrations (rank INTEGER NOT NULL UNIQUE, NAME TEXT NOT NULL);")
	row := db.QueryRow("SELECT rank FROM __migrations ORDER BY rank DESC LIMIT 1;")
	var rank int

	err := row.Scan(&rank)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	pattern := filepath.Join(srcPath, "v_*__*.sql")
	matches, err := filepath.Glob(pattern)

	if err != nil {
		return err
	}

	idx := 0
	var ver int

	for i, match := range matches {
		match = filepath.Base(match)
		_ver := match[2:strings.Index(match, "__")]
		ver, err = strconv.Atoi(_ver)

		if err != nil {
			return errors.New("bad migration filename pattern. version " + _ver)
		}

		if ver > rank {
			idx = i
			break
		}

		idx = i + 1
	}

	if len(matches[idx:]) == 0 {
		fmt.Println("migrate: no migrations found, skipping.")
		return nil
	}

	if err != nil {
		return err
	}

	for _, match := range matches[idx:] {
		ver, err := version(match)
		if err != nil {
			return errors.New("bad migration filename. " + filepath.Base(match))
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		schema, err := readSchema(filepath.Join(srcPath, filepath.Base(match)))
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(schema)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO __migrations (rank, name) values (?, ?)", ver, name(match))
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func version(m string) (int, error) {
	_m := filepath.Base(m)
	ver := _m[2:strings.Index(_m, "__")]
	return strconv.Atoi(ver)
}
func name(m string) string {
	_m := filepath.Base(m)
	return _m[strings.Index(_m, "__")+2 : strings.Index(_m, ".sql")]
}

func readSchema(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	schema, err := io.ReadAll(file)

	return string(schema), nil
}
