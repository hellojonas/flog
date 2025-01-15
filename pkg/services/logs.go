package services

import (
	"database/sql"
	"path/filepath"
	"time"
)

type Log struct {
	Name      string
	AppId     int64
	CreatedAt time.Time
}

type LogCreateInput struct {
	Name  string
	AppId int64
}

type LogService struct {
	logDir string
	db     *sql.DB
}

func (ls *LogService) ResolveLogFile(appName string, logName string) string {
	return filepath.Join(ls.logDir, appName, logName)
}

func NewLogService(db *sql.DB, logDir string) *LogService {
	return &LogService{
		db:     db,
		logDir: logDir,
	}
}

func (ls *LogService) CreateLog(data LogCreateInput) error {
	_, err := ls.db.Exec("INSERT INTO logs (name, application_id) VALUES (?, ?);", data.Name, data.AppId)

	if err != nil {
		return err
	}

	return nil
}

func (ls *LogService) ListAppLogs(appId int64) ([]Log, error) {
	rows, err := ls.db.Query("SELECT name, created_at FROM logs WHERE application_id = ?;", appId)

	if err != nil {
		return nil, err
	}

	var name string
	var createdAt time.Time

	var logs []Log
	for rows.Next() {
		err := rows.Scan(&name, &createdAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, Log{
			Name:      name,
			AppId:     appId,
			CreatedAt: createdAt,
		})
	}

	return logs, nil
}
