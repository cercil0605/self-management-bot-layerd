package service

import (
	"fmt"
	"self-management-bot/repository"
)

func AddTaskService(userID, title string) error {
	return repository.AddTask(userID, title)
}
func GetTaskService(userID string) ([]repository.Task, error) {
	return repository.FindTaskByUserID(userID)
}
func CompleteTaskService(userID string, DoneTaskNumber int) error {
	tasks, err := GetTaskService(userID)
	// 内部エラー
	if err != nil {
		return fmt.Errorf("タスク取得に失敗: %w", err)
	}
	if len(tasks) == 0 {
		return fmt.Errorf("タスクが1件も登録されていません")
	}
	// タスク存在
	if DoneTaskNumber < 0 || DoneTaskNumber >= len(tasks) {
		return fmt.Errorf("指定されたタスク番号は存在しません")
	}
	return repository.CompleteTask(tasks[DoneTaskNumber].ID)
}
func DeleteTaskService(userID string, DeleteTaskNumber int) error {
	tasks, err := GetTaskService(userID)
	// 内部エラー
	if err != nil {
		return fmt.Errorf("タスク取得に失敗: %w", err)
	}
	if len(tasks) == 0 {
		return fmt.Errorf("タスクが1件も登録されていません")
	}
	// タスク存在
	if DeleteTaskNumber < 0 || DeleteTaskNumber >= len(tasks) {
		return fmt.Errorf("指定されたタスク番号は存在しません")
	}
	return repository.DeleteTask(tasks[DeleteTaskNumber].ID)
}
