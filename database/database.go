package database

import (
	"com/parser/parser"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

var DB *sql.DB

func init() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	handleDbError(err)
	DB = db
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS logs (id SERIAL)")
	handleDbError(err)
	_, err = DB.Exec("ALTER TABLE logs ADD COLUMN IF NOT EXISTS timestamp_start TIMESTAMP")
	handleDbError(err)
	_, err = DB.Exec("ALTER TABLE logs ADD COLUMN IF NOT EXISTS timestamp_end TIMESTAMP")
	handleDbError(err)
}

type LogEntry struct {
	Id int
	TimestampStart pq.NullTime
	TimestampFinish pq.NullTime
}

func InsertLog(log parser.Log) {
	sql := "INSERT INTO logs (id, timestamp_start, timestamp_end) VALUES ($1, $2, $3)"
	_, err := DB.Exec(sql, log.Id, log.TimestampStart, log.TimestampFinish)
	handleDbError(err)
}

func FetchLogs() ([]parser.Log, error) {
	logs := make([]parser.Log, 0)

	rows, err := DB.Query("SELECT id, timestamp_start, timestamp_end FROM logs")
	if err != nil {
		return logs, err
	}
	defer rows.Close()

	for rows.Next() {
		var log parser.Log
		err := rows.Scan(&log.Id, &log.TimestampStart, &log.TimestampFinish)
		if err != nil {
			return logs, err
		}
		handleDbError(err)
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return logs, err
	}

	return logs, nil
}

func handleDbError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err)
		return true
	}
	return false
}
