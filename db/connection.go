package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
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
	// マイグレーション
	sql, err := ioutil.ReadFile("db/migration.sql")
	if err != nil {
		log.Println("⚠️ マイグレーションSQL読み込み失敗:", err)
		return err
	}
	_, err = DB.Exec(string(sql))
	if err != nil {
		log.Println("❌ マイグレーション実行失敗:", err)
		return err
	}
	return DB.Ping()
}
