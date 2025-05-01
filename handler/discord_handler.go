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
	case strings.HasPrefix(content, "!chat"):
		HandleChat(s, m, content)
	case strings.HasPrefix(content, "!reset"):
		HandleReset(s, m)
	case strings.HasPrefix(content, "!confirm reset"):
		HandleConfirm(s, m)
	}
}

func HandleAdd(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	title := strings.TrimPrefix(content, "!add ")
	if len(title) == 0 {
		replyToUser(s, m.ChannelID, m.Author.ID, "```⚠️ タスク内容を追加してください```")
		return
	}
	err := service.AddTaskService(m.Author.ID, title)
	if err != nil {
		replyToUser(s, m.ChannelID, m.Author.ID, "```❌ タスク登録失敗```")
		return
	}
	replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```⭕️ タスク追加: %s```", title))
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
	for i, task := range tasks {
		msg.WriteString(fmt.Sprintf("⌛️ [%02d] %s\n", i, task.Title))
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
	var msg strings.Builder
	msg.WriteString("```✅ タスク完了！お疲れ様です！\n")
	if len(tasks) == 0 {
		msg.WriteString("\n🎉 もう残ってるタスクはありません！今日もよく頑張った！```")
	} else {
		msg.WriteString("\n📝 残りのタスク:\n")
		for i, task := range tasks {
			msg.WriteString(fmt.Sprintf("⌛️ [%02d] %s\n", i, task.Title))
		}
		msg.WriteString("```")
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
	s.ChannelTyping(m.ChannelID)
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
