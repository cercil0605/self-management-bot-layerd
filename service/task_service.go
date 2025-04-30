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
	if err != nil {
		return fmt.Errorf("タスク取得に失敗: %w", err)
	}
	if len(tasks) == 0 {
		return fmt.Errorf("タスクが1件も登録されていません")
	}
	if DoneTaskNumber < 0 || DoneTaskNumber >= len(tasks) {
		return fmt.Errorf("指定されたタスク番号は存在しません")
	}
	return repository.CompleteTask(tasks[DoneTaskNumber].ID)
}
