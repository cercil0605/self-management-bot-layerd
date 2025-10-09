package service

// タスク関連のCRUD処理
import (
	"fmt"
	"self-management-bot/client"
	"self-management-bot/lib/errors"
	"self-management-bot/repository"
	"strings"
)

func AddTaskService(userID, title string, priorityID int) error {
	err := repository.AddTask(userID, title, priorityID)
	if err != nil {
		return errors.NewAppError("AddTaskService", err)
	}
	return nil
}

// GetTaskService 今日のタスクを取得
func GetTaskService(userID string) ([]repository.Task, error) {
	tasks, err := repository.FindTaskByUserID(userID, "today")
	if err != nil {
		return nil, errors.NewAppError("GetTaskService", err)
	}
	return tasks, nil
}

func GetYesterdayTaskService(userID string) ([]repository.Task, error) {
	tasks, err := repository.FindTaskByUserID(userID, "yesterday")
	if err != nil {
		return nil, errors.NewAppError("GetYesterdayTaskService", err)
	}
	return tasks, nil
}

func UpdateTaskService(userID string, TaskNumber int, title string, priorityID *int) error {
	tasks, err := GetTaskService(userID)
	// 内部エラー
	if err != nil {
		return errors.NewAppError("UpdateTaskService", err)
	}
	if len(tasks) == 0 {
		return errors.NewAppError("UpdateTaskService", fmt.Errorf("タスクが1件も登録されていません"))
	}
	if TaskNumber < 0 || TaskNumber >= len(tasks) {
		return errors.NewAppError("UpdateTaskService", fmt.Errorf("指定されたタスク番号は存在しません"))
	}
	err = repository.UpdateTask(tasks[TaskNumber].ID, title, priorityID)
	if err != nil {
		return errors.NewAppError("UpdateTaskService", err)
	}
	return nil
}

func CompleteTaskService(userID string, DoneTaskNumber int) error {
	tasks, err := GetTaskService(userID)
	// 内部エラー
	if err != nil {
		return errors.NewAppError("CompleteTaskService", fmt.Errorf("タスク取得に失敗: %w", err))
	}
	if len(tasks) == 0 {
		return errors.NewAppError("CompleteTaskService", fmt.Errorf("タスクが1件も登録されていません"))
	}
	// タスク存在
	if DoneTaskNumber < 0 || DoneTaskNumber >= len(tasks) {
		return errors.NewAppError("CompleteTaskService", fmt.Errorf("指定されたタスク番号は存在しません"))
	}
	if err := repository.CompleteTask(tasks[DoneTaskNumber].ID); err != nil {
		return errors.NewAppError("CompleteTaskService", err)
	}
	return nil
}

func DeleteTaskService(userID string, DeleteTaskNumber int) error {
	tasks, err := GetTaskService(userID)
	// 内部エラー
	if err != nil {
		return errors.NewAppError("DeleteTaskService", fmt.Errorf("タスク取得に失敗: %w", err))
	}
	if len(tasks) == 0 {
		return errors.NewAppError("DeleteTaskService", fmt.Errorf("タスクが1件も登録されていません"))
	}
	// タスク存在
	if DeleteTaskNumber < 0 || DeleteTaskNumber >= len(tasks) {
		return errors.NewAppError("DeleteTaskService", fmt.Errorf("指定されたタスク番号は存在しません"))
	}
	if err := repository.DeleteTask(tasks[DeleteTaskNumber].ID); err != nil {
		return errors.NewAppError("DeleteTaskService", err)
	}
	return nil
}

// ChatWithContext 今日のタスク状況について
func ChatWithContext(userID, input string) (string, error) {
	pending, err := repository.FindPendingTaskByUser(userID)
	if err != nil {
		return "❌ ユーザーのタスク取得に失敗しました(Pending)", errors.NewAppError("ChatWithContext", err)
	}
	completed, err := repository.FindCompletedTodayTaskByUser(userID)
	if err != nil {
		return "❌ ユーザーのタスク取得に失敗しました(Completed)", errors.NewAppError("ChatWithContext", err)
	}
	prompt := CreateChatPrompt(pending, completed, input)
	res, err := client.GetGeminiResponse(prompt)
	if err != nil {
		return "❌ 応答に失敗しました(LLM)", errors.NewAppError("ChatWithContext", err)
	}
	return res, nil
}

// CreateChatPrompt 今日の完了状況をプロンプト化する
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
func ResetTodayTasks(userID string) (int, error) {
	n, err := repository.DeleteTodayTasks(userID)
	if err != nil {
		return 0, errors.NewAppError("ResetTodayTasks", err)
	}
	return n, nil
}
func ResetAllTasks(userID string) (int, error) {
	n, err := repository.DeleteAllTasksByUser(userID)
	if err != nil {
		return 0, errors.NewAppError("ResetAllTasks", err)
	}
	return n, nil
}

type ReminderMessage struct {
	Content string // LLMからのメッセージ
	UserID  string // ユーザID
}

// FixedTimeReminder 定期リマインダ送信
func FixedTimeReminder() ([]ReminderMessage, error) {
	userInfo, err := repository.FindAllUser()
	if err != nil {
		fmt.Println("❌ ユーザ情報取得失敗:", err)
		return nil, errors.NewAppError("FixedTimeReminder", err)
	}
	// TODO 登録したすべてのユーザに送信するようにする
	tasks, err := GetYesterdayTaskService(userInfo[0])
	if err != nil {
		fmt.Printf("❌ タスク取得失敗 userID=%s: %v\n", userInfo[0], err)
		return nil, errors.NewAppError("FixedTimeReminder", err)
	}
	var prompt strings.Builder
	// プロンプト
	prompt.WriteString("あなたは自己管理を支援するプロフェッショナルなコーチです。\n")
	prompt.WriteString("昨日のタスクの実行状況をふまえ、今日を気持ちよくスタートできるように前向きで実用的なアドバイスを与えてください。\n")
	prompt.WriteString("以下のルールに従ってください：\n")
	prompt.WriteString("- 昨日の達成を簡潔に肯定的に振り返る（完了したタスクがあれば）\n")
	prompt.WriteString("- 昨日未完了だったタスクがあれば、それをどう今日活かすか助言する\n")
	prompt.WriteString("- アドバイスは1〜3個、シンプルかつ実行可能なものにする\n\n")

	// 昨日のタスクの整理
	prompt.WriteString("【昨日のタスク状況】\n")

	hasCompleted := false
	hasPending := false
	// TODO Goに任せるんじゃなくてDB操作に任せたい（あとで）
	for _, task := range tasks {
		switch task.Status {
		case "completed":
			if !hasCompleted {
				prompt.WriteString("▼完了したタスク：\n")
				hasCompleted = true
			}
			prompt.WriteString("- " + task.Title + "\n")
		case "pending":
			if !hasPending {
				prompt.WriteString("▼未完了のタスク：\n")
				hasPending = true
			}
			prompt.WriteString("- " + task.Title + "\n")
		}
	}
	// 何もタスクをこなしてない時
	if !hasCompleted {
		prompt.WriteString("▼完了したタスク：\n（完了したタスクはありません）\n")
	}
	if !hasPending {
		prompt.WriteString("▼未完了のタスク：\n（未完了のタスクはありません）\n")
	}
	prompt.WriteString("\nこの情報をふまえて、今日をポジティブに始めるためのメッセージを作成してください。\n")
	res, err := client.GetGeminiResponse(prompt.String())
	if err != nil {
		fmt.Printf("❌ LLM応答失敗 userID=%s: %v\n", userInfo[0], err)
		return nil, errors.NewAppError("FixedTimeReminder", err)
	}

	fmt.Println("✅ リマインド生成成功")

	msg := ReminderMessage{
		Content: res,
		UserID:  userInfo[0],
	}
	return []ReminderMessage{msg}, nil
}
