package repository

import (
	"fmt"
	"self-management-bot/db"
	"self-management-bot/lib/errors"
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
	query := `INSERT INTO tasks (user_id, title, priority_id, status) VALUES ($1, $2, $3, 'pending')`
	_, err := db.DB.Exec(query, userID, title, priorityID)
	if err != nil {
		return errors.NewAppError("AddTask", err)
	}
	return nil
}

// FindTaskByUserID 完了状況問わずタスクを出力
func FindTaskByUserID(userID string, when string) ([]Task, error) {
	baseQuery := `
		SELECT id, title, status, priority_id FROM tasks
		WHERE user_id = $1 %s
		ORDER BY
			CASE status
				WHEN 'pending' THEN 0
				WHEN 'completed' THEN 1
			END,
			priority_id ASC`
	// 日付に応じてSQL文を変えて絞り込む
	var dateCondition string
	if when == "today" {
		// 未完了タスクは日付問わず表示する
		dateCondition = "AND (status = 'pending' OR (status = 'completed' AND created_at::date = CURRENT_DATE))"
	} else if when == "yesterday" {
		dateCondition = "AND created_at >= CURRENT_DATE - INTERVAL '1 day' AND created_at < CURRENT_DATE"
	} else {
		dateCondition = ""
	}
	query := fmt.Sprintf(baseQuery, dateCondition)

	var tasks []Task
	err := db.DB.Select(&tasks, query, userID)
	if err != nil {
		return nil, errors.NewAppError("FindTaskByUserID", err)
	}
	return tasks, nil
}
func UpdateTask(taskID int, title string, priorityID *int) error {
	var query string
	var args []interface{}
	// 1 value
	if priorityID == nil {
		query = `UPDATE tasks SET title = $1 WHERE id = $2`
		args = []interface{}{title, taskID}
	} else { // 2 value
		if title == "" {
			query = `UPDATE tasks SET priority_id = $1 WHERE id = $2`
			args = []interface{}{*priorityID, taskID}
		} else {
			query = `UPDATE tasks SET title = $1, priority_id = $2 WHERE id = $3`
			args = []interface{}{title, *priorityID, taskID}
		}
	}
	_, err := db.DB.Exec(query, args...)
	if err != nil {
		errors.NewAppError("UpdateTask", err)
	}
	return nil
}
func CompleteTask(taskID int) error {
	query := `UPDATE tasks SET status = 'completed' WHERE id = $1`
	_, err := db.DB.Exec(query, taskID)
	if err != nil {
		return errors.NewAppError("CompleteTask", err)
	}
	return nil
}
func DeleteTask(taskID int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := db.DB.Exec(query, taskID)
	if err != nil {
		return errors.NewAppError("DeleteTask", err)
	}
	return nil
}

// FindCompletedTodayTaskByUser 今日の完了済みタスク
func FindCompletedTodayTaskByUser(userID string) ([]Task, error) {
	query := `SELECT id, title, status FROM tasks 
			  WHERE user_id = $1 AND status = 'completed' AND created_at::date = CURRENT_DATE
			  ORDER BY created_at`
	var tasks []Task
	err := db.DB.Select(&tasks, query, userID)
	if err != nil {
		return nil, errors.NewAppError("FindCompletedTodayTaskByUser", err)
	}
	return tasks, nil
}

// FindPendingTaskByUser 待ちタスク
func FindPendingTaskByUser(userID string) ([]Task, error) {
	query := `SELECT id, title, status FROM tasks 
			  WHERE user_id = $1 AND status = 'pending'
			  ORDER BY created_at`
	var tasks []Task
	err := db.DB.Select(&tasks, query, userID)
	if err != nil {
		return nil, errors.NewAppError("FindPendingTaskByUser", err)
	}
	return tasks, nil
}

// FindAllUser ユーザIDを全て探す
func FindAllUser() ([]string, error) {
	query := `SELECT DISTINCT user_id FROM tasks`
	var userIDs []string
	err := db.DB.Select(&userIDs, query)
	if err != nil {
		return nil, errors.NewAppError("FindAllUser", err)
	}
	return userIDs, nil
}

func DeleteTodayTasks(userID string) (int, error) {
	query := `
		DELETE FROM tasks
		WHERE user_id = $1 AND created_at::date = CURRENT_DATE
	`
	res, err := db.DB.Exec(query, userID)
	if err != nil {
		return 0, errors.NewAppError("DeleteTodayTasks", err)
	}
	rows, rerr := res.RowsAffected()
	if rerr != nil {
		return 0, errors.NewAppError("DeleteTodayTasks", rerr)
	}
	return int(rows), nil
}

func DeleteAllTasksByUser(userID string) (int, error) {
	query := `DELETE FROM tasks WHERE user_id = $1`
	res, err := db.DB.Exec(query, userID)
	if err != nil {
		return 0, errors.NewAppError("DeleteAllTasksByUser", err)
	}
	rows, rerr := res.RowsAffected()
	if rerr != nil {
		return 0, errors.NewAppError("DeleteAllTasksByUser", rerr)
	}
	return int(rows), nil
}
