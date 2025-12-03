package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"student-management/config"
)

var GradeDB *sql.DB

func InitGradeDB(cfg *config.Config) error {
	var err error
	GradeDB, err = sql.Open("postgres", cfg.GetGradeDBURL())
	if err != nil {
		return err
	}
	return GradeDB.Ping()
}
