package service

import "self-management-bot/repository"

func AddTask(userID, title string) error {
	return repository.AddTask(userID, title)
}
