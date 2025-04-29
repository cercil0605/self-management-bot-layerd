package service

import "self-management-bot/repository"

func AddTaskService(userID, title string) error {
	return repository.AddTask(userID, title)
}
func GetTaskService(userID string) ([]repository.Task, error) {
	return repository.FindTaskByUserID(userID)
}
