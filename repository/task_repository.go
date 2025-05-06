package repository

import (
	"fmt"
	"self-management-bot/db"
)

type Task struct {
	ID         int    `db:"id"`
	UserID     string `db:"user_id"`
	Title      string `db:"title"`
	PriorityID int    `db:"priority_id"`
	Status     string `db:"status"`
}

type Priority struct {
	ID    int    `db:"id"`
	Code  string `db:"code"`
	Emoji string `db:"emoji"`
}

func AddTask(userID, title string, priorityID int) error {
	query := `INSERT INTO tasks (user_id, title, priority_id,status) VALUES ($1, $2, $3,'pending')`
	_, err := db.DB.Exec(query, userID, title, priorityID)
	if err != nil {
		fmt.Println("❌ AddTask error:", err)
	}
	return err
}

// FindTaskByUserID 完了状況問わず今日のタスクを表示
func FindTaskByUserID(userID string) ([]Task, error) {
	query := `SELECT id,title,status,priority_id FROM tasks 
                       WHERE user_id = $1  AND created_at::date = CURRENT_DATE
                       ORDER BY 
                           CASE status
                           WHEN 'pending' THEN 0
                           WHEN 'completed' THEN 1
					   END,
					   priority_id ASC`
	var tasks []Task
	err := db.DB.Select(&tasks, query, userID)
	return tasks, err
}
func UpdateTask(taskID int, title string, priorityID int) error {
	query := `UPDATE tasks SET title = $1 , priority_id = $2 WHERE id = $3`
	_, err := db.DB.Exec(query, title, priorityID, taskID)
	return err
}
func CompleteTask(taskID int) error {
	query := `UPDATE tasks SET status = 'completed' WHERE id = $1`
	_, err := db.DB.Exec(query, taskID)
	return err
}
func DeleteTask(taskID int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := db.DB.Exec(query, taskID)
	return err
}

// FindCompletedTodayTaskByUser 今日の完了済みタスク
func FindCompletedTodayTaskByUser(userID string) ([]Task, error) {
	query := `SELECT id,title,status FROM tasks 
                       WHERE user_id = $1 AND status = 'completed' AND created_at::date = CURRENT_DATE
                       ORDER BY created_at `
	var tasks []Task
	err := db.DB.Select(&tasks, query, userID)
	return tasks, err
}

// FindPendingTodayTaskByUser 今日の待ちタスク
func FindPendingTodayTaskByUser(userID string) ([]Task, error) {
	query := `SELECT id,title,status FROM tasks 
                       WHERE user_id = $1 AND status = 'pending' AND created_at::date = CURRENT_DATE
                       ORDER BY created_at `
	var tasks []Task
	err := db.DB.Select(&tasks, query, userID)
	return tasks, err
}

func DeleteTodayTasks(userID string) (int, error) {
	query := `
		DELETE FROM tasks
		WHERE user_id = $1 AND created_at::date = CURRENT_DATE
	`
	res, err := db.DB.Exec(query, userID)
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()
	return int(rows), nil
}

func DeleteAllTasksByUser(userID string) (int, error) {
	query := `DELETE FROM tasks WHERE user_id = $1`
	res, err := db.DB.Exec(query, userID)
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()
	return int(rows), nil
}
