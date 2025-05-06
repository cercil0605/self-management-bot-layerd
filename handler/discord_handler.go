package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"self-management-bot/service"
	"strconv"
	"strings"
	"time"
)

var resetAllConfirm = make(map[string]time.Time)

// 優先度チェック
var priorityMap = map[string]int{
	"P1": 1,
	"P2": 2,
	"P3": 3,
	"P4": 4,
}
var priorityEmoji = map[int]string{
	1: "🔴", // P1
	2: "🟡", // P2
	3: "🟢", // P3
	4: "🔵", // P4
}

func replyToUser(s *discordgo.Session, chID, userID, message string) {
	_, err := s.ChannelMessageSend(chID, fmt.Sprintf("<@%s>\n%s", userID, message))
	if err != nil {
		fmt.Printf("⚠️ Discord送信エラー: %v\n", err)
	}
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := strings.TrimSpace(m.Content)

	switch {
	case strings.HasPrefix(content, "!add "):
		HandleAdd(s, m, content)
	case strings.HasPrefix(content, "!list"):
		HandleList(s, m)
	case strings.HasPrefix(content, "!done "):
		HandleComplete(s, m, content)
	case strings.HasPrefix(content, "!delete"):
		HandleDelete(s, m, content)
	case strings.HasPrefix(content, "!chat "):
		HandleChat(s, m, content)
	case strings.HasPrefix(content, "!reset"):
		HandleReset(s, m)
	case strings.HasPrefix(content, "!confirm reset"):
		HandleConfirm(s, m)
	case strings.HasPrefix(content, "!edit "):
		HandleEdit(s, m, content)
	case strings.HasPrefix(content, "!help"):
		HandleHelp(s, m)
	}
}

func HandleAdd(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	args := strings.Fields(strings.TrimPrefix(content, "!add"))
	if len(args) == 0 {
		replyToUser(s, m.ChannelID, m.Author.ID, "```⚠️ タスク内容を追加してください```")
		return
	}
	// 優先度を表す部分だけTrim
	priorityID := 4 // default
	priorityInput := strings.ToUpper(args[len(args)-1])
	if pid, ok := priorityMap[priorityInput]; ok {
		priorityID = pid
		args = args[:len(args)-1]
	}
	title := strings.Join(args, " ")
	err := service.AddTaskService(m.Author.ID, title, priorityID)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, "```❌ タスク登録失敗```")
		return
	}
	replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```⭕️ タスク追加: %s 優先度： %d (%s)```", title, priorityID, priorityEmoji[priorityID]))
}

func HandleList(s *discordgo.Session, m *discordgo.MessageCreate) {
	tasks, err := service.GetTaskService(m.Author.ID)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, "```❌ タスク取得失敗```")
		return
	}
	if len(tasks) == 0 {
		replyToUser(s, m.ChannelID, m.Author.ID, "```📭 タスクが登録されていません```")
		return
	}
	var msg strings.Builder
	msg.WriteString("今日のTodoです！\n```")
	completedFlag := false
	for i, task := range tasks {
		if task.Status == "pending" {
			if i == 0 {
				msg.WriteString(fmt.Sprintf("📝 未完了のタスク\n"))
			}
			msg.WriteString(fmt.Sprintf("%s ⌛️ [%02d] %s\n", priorityEmoji[task.PriorityID], i, task.Title))
		} else if task.Status == "completed" {
			if completedFlag == false {
				msg.WriteString(fmt.Sprintf("\n✅ 完了済みのタスク\n"))
				completedFlag = true
			}
			msg.WriteString(fmt.Sprintf("✅ [%02d] %s\n", i, task.Title))
		}
	}
	msg.WriteString("```")
	replyToUser(s, m.ChannelID, m.Author.ID, msg.String())
}

func HandleComplete(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	arg := strings.TrimPrefix(content, "!done ")
	DoneTaskNumber, err := strconv.Atoi(arg)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, "```❌ 数字を指定してください```")
		return
	}
	err = service.CompleteTaskService(m.Author.ID, DoneTaskNumber)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ %s```", err.Error()))
		return
	}
	tasks, err := service.GetTaskService(m.Author.ID)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, "```✅ タスク完了！\n⚠️ 残りのタスク取得に失敗しました```")
		return
	}
	// 内容出力
	var msg strings.Builder
	msg.WriteString("```✅ タスク完了！お疲れ様です！\n")
	hasPending := false
	for i, task := range tasks {
		if task.Status == "pending" {
			if !hasPending {
				msg.WriteString("\n📝 残りのタスク:\n")
				hasPending = true
			}
			msg.WriteString(fmt.Sprintf("⌛️ [%02d] %s\n", i, task.Title))
		}
	}
	if hasPending {
		msg.WriteString("```")
	} else {
		msg.WriteString("\n🎉 もう残ってるタスクはありません！今日もよく頑張った！```")
	}
	replyToUser(s, m.ChannelID, m.Author.ID, msg.String())
}

func HandleDelete(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	arg := strings.TrimPrefix(content, "!delete ")
	DeleteNumber, err := strconv.Atoi(arg)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, "```❌ 数字を指定してください```")
		return
	}
	err = service.DeleteTaskService(m.Author.ID, DeleteNumber)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ %s```", err.Error()))
		return
	}
	replyToUser(s, m.ChannelID, m.Author.ID, "```⭕️ タスク削除しました```")
}

func HandleChat(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	arg := strings.TrimPrefix(content, "!chat ")
	if len(strings.TrimSpace(arg)) == 0 {
		replyToUser(s, m.ChannelID, m.Author.ID, "```❌ メッセージを入力してください```")
		return
	}
	err := s.ChannelTyping(m.ChannelID)
	if err != nil {
		return
	}
	reply, err := service.ChatWithContext(m.Author.ID, arg)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ %s```", err.Error()))
		return
	}
	replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```\n%s\n```", reply))
}

func HandleReset(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, "!reset all") {
		resetAllConfirm[m.Author.ID] = time.Now().Add(10 * time.Minute)
		replyToUser(s, m.ChannelID, m.Author.ID,
			"```⚠️ 本当に全タスク（過去含む）を削除しますか？\n削除するには '!confirm reset' と入力してください。（10分以内）```")
		return
	}
	count, err := service.ResetTodayTasks(m.Author.ID)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ 今日のリセット失敗: %s```", err.Error()))
		return
	}
	replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```✅ 今日のタスクを %d 件削除しました```", count))
}

func HandleConfirm(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	expiry, ok := resetAllConfirm[userID]
	if !ok || time.Now().After(expiry) {
		delete(resetAllConfirm, userID)
		replyToUser(s, m.ChannelID, userID, "```⚠️ '!reset all' の確認時間が切れました。再度実行してください。```")
		return
	}

	count, err := service.ResetAllTasks(userID)
	if err != nil {
		replyToUser(s, m.ChannelID, userID, fmt.Sprintf("```❌ 全削除に失敗しました: %s```", err.Error()))
		return
	}
	delete(resetAllConfirm, userID)
	replyToUser(s, m.ChannelID, userID, fmt.Sprintf("```✅ 全タスクを %d 件削除しました```", count))
}

func HandleEdit(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	arg := strings.TrimPrefix(content, "!edit ")
	fields := strings.Fields(arg)
	if len(fields) < 2 {
		replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```⚠️ コマンドの形式が正しくありません。\n例: `!edit 1 新しい内容` ```"))
		return
	}
	IndexNumber, err := strconv.Atoi(fields[0])
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, "```❌ 数字を指定してください```")
		return
	}
	newTitle := fields[1]
	// 優先度の値を設定
	var newPriority int
	if fields[2] == "" {
		newPriority = 4 // default
	} else {
		newPriority = priorityMap[fields[2]] // 代入
	}
	err = service.UpdateTaskService(m.Author.ID, IndexNumber, newTitle, newPriority)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ タスクの編集に失敗しました: %s```", err.Error()))
		return
	}
	replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```✅ 指定されたToDoを編集しました```"))
}

func HandleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpText := "# 入れてくれてありがとう！\n 💻 コマンド一覧だよ\n" + "```" +
		"!add <タスク名> [P1~P4]    : タスクを追加（例: !add 宿題 P1）\n" +
		"!list                      : 今日のタスクを一覧表示\n" +
		"!done <番号>              : タスクを完了扱いに\n" +
		"!delete <番号>           : タスクを削除\n" +
		"!reset                    : 今日のタスクを全削除\n" +
		"!reset all               : 全タスクを削除（確認付き）\n" +
		"!confirm reset           : 全削除を確定\n" +
		"!chat <メッセージ>        : AIと会話\n" +
		"!help                     : このヘルプを表示\n" +
		"```"
	replyToUser(s, m.ChannelID, m.Author.ID, helpText)
}
