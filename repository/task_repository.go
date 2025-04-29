package repository

import (
	"fmt"
	"self-management-bot/db"
)

type Task struct {
	ID     int    `db:"id"`
	UserID string `db:"user_id"`
	Title  string `db:"title"`
	Status string `db:"status"`
}

func AddTask(userID, title string) error {
	query := `INSERT INTO tasks (user_id, title, status) VALUES ($1, $2, 'pending')`
	_, err := db.DB.Exec(query, userID, title)
	if err != nil {
		fmt.Println("❌ AddTask error:", err)
	}
	return err
}
