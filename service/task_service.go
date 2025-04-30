package service

import (
	"fmt"
	"self-management-bot/client"
	"self-management-bot/repository"
	"strings"
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

func ChatWithContext(userID, input string) (string, error) {
	pending, err := repository.FindTaskByUserID(userID)
	if err != nil {
		return "❌ ユーザーのタスク取得に失敗しました(Pending)", err
	}
	completed, err := repository.FindCompletedTaskByUser(userID)
	if err != nil {
		return "❌ ユーザーのタスク取得に失敗しました(Completed)", err
	}
	prompt := CreateChatPrompt(pending, completed, input)
	res, err := client.GetLLMResponse(prompt)
	if err != nil {
		return "❌ 応答に失敗しました(LLM)", err
	}
	return res, nil
}

func CreateChatPrompt(pending []repository.Task, completed []repository.Task, input string) string {
	var prompt strings.Builder
	prompt.WriteString("あなたは，自己管理を支援するメンズコーチです．\n\n")
	prompt.WriteString("【未完了のタスク】\n")
	if len(pending) == 0 {
		prompt.WriteString("（未完了のタスクはありません）\n")
	} else {
		for _, t := range pending {
			prompt.WriteString("- " + t.Title + "\n")
		}
	}
	prompt.WriteString("\n【最近完了したタスク】\n")
	if len(completed) == 0 {
		prompt.WriteString("（完了したタスクはありません）\n")
	} else {
		for _, t := range completed {
			prompt.WriteString("- " + t.Title + "\n")
		}
	}
	prompt.WriteString("\n【ユーザーの質問】\n")
	prompt.WriteString(input + "\n")
	prompt.WriteString("\n上記を踏まえてアドバイスせよ．")
	return prompt.String()
}
