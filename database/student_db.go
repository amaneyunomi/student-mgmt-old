package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"student-management/config"
)

var StudentDB *sql.DB

func InitStudentDB(cfg *config.Config) error {
	var err error
	StudentDB, err = sql.Open("postgres", cfg.GetStudentDBURL())
	if err != nil {
		return err
	}
	return StudentDB.Ping()
}
