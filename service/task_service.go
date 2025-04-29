package service

import "self-management-bot/repository"

func AddTaskService(userID, title string) error {
	return repository.AddTask(userID, title)
}
