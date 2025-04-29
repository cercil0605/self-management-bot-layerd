package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

var DB *sqlx.DB

// Init connection
func Init() error {
	var err error
	DB, err = sqlx.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	return DB.Ping()
}
