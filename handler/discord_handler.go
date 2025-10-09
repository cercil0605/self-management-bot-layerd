package handler

import (
	"fmt"
	"log"
	"self-management-bot/lib/errors"
	"self-management-bot/service"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const COMMAND_PREFIX = "!"

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

func replyToUser(s *discordgo.Session, chID, userID, message string) error {
	_, err := s.ChannelMessageSend(chID, fmt.Sprintf("<@%s>\n%s", userID, message))
	if err != nil {
		return errors.NewAppError("replyToUser.ChannelMessageSend", err)
	}
	return nil
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := messageCreate(s, m); err != nil {
		err = errors.NewAppError("MessageCreate", err)
		log.Printf("❌ %s\n", err.Error())
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.Author.ID == s.State.User.ID {
		return nil
	}

	content := strings.TrimSpace(m.ContentWithMentionsReplaced())

	switch {
	case strings.HasPrefix(content, fmt.Sprintf("!add ", COMMAND_PREFIX)):
		if err := HandleAdd(s, m, content); err != nil {
			return errors.NewAppError("HandleAdd", err)
		}
	case strings.HasPrefix(content, fmt.Sprintf("!list", COMMAND_PREFIX)):
		if err := HandleList(s, m); err != nil {
			return errors.NewAppError("HandleList", err)
		}
	case strings.HasPrefix(content, fmt.Sprintf("!done ", COMMAND_PREFIX)):
		if err := HandleComplete(s, m, content); err != nil {
			return errors.NewAppError("HandleComplete", err)
		}
	case strings.HasPrefix(content, fmt.Sprintf("!delete", COMMAND_PREFIX)):
		if err := HandleDelete(s, m, content); err != nil {
			return errors.NewAppError("HandleDelete", err)
		}
	case strings.HasPrefix(content, fmt.Sprintf("!chat ", COMMAND_PREFIX)):
		if err := HandleChat(s, m, content); err != nil {
			return errors.NewAppError("HandleChat", err)
		}
	case strings.HasPrefix(content, fmt.Sprintf("!reset", COMMAND_PREFIX)):
		if err := HandleReset(s, m); err != nil {
			return errors.NewAppError("HandleReset", err)
		}
	case strings.HasPrefix(content, fmt.Sprintf("!confirm reset", COMMAND_PREFIX)):
		if err := HandleConfirm(s, m); err != nil {
			return errors.NewAppError("HandleConfirm", err)
		}
	case strings.HasPrefix(content, fmt.Sprintf("!edit ", COMMAND_PREFIX)):
		if err := HandleEdit(s, m, content); err != nil {
			return errors.NewAppError("HandleEdit", err)
		}
	case strings.HasPrefix(content, fmt.Sprintf("!help", COMMAND_PREFIX)):
		if err := HandleHelp(s, m); err != nil {
			return errors.NewAppError("HandleHelp", err)
		}
	}
	return nil
}

func HandleAdd(s *discordgo.Session, m *discordgo.MessageCreate, content string) error {
	args := strings.Fields(strings.TrimPrefix(content, fmt.Sprintf("%sadd", COMMAND_PREFIX)))
	if len(args) == 0 {
		if err := replyToUser(s, m.ChannelID, m.Author.ID, "```⚠️ タスク内容を追加してください```"); err != nil {
			return errors.NewAppError("noTaskContent", err)
		}
		return errors.NewAppError("noTaskContent", fmt.Errorf("no task content"))
	}

	// 優先度を表す部分だけTrim
	priorityID := 4 // default
	priorityInput := strings.ToUpper(args[len(args)-1])
	if pid, ok := priorityMap[priorityInput]; ok {
		priorityID = pid
		args = args[:len(args)-1]
	}
	title := strings.Join(args, " ")

	if err := service.AddTaskService(m.Author.ID, title, priorityID); err != nil {
		// ユーザーへの通知が失敗した場合も別名で返す
		if err := replyToUser(s, m.ChannelID, m.Author.ID, "```❌ タスク登録失敗```"); err != nil {
			return errors.NewAppError("failedToAddTask", err)
		}
		return err
	}

	if err := replyToUser(
		s, m.ChannelID, m.Author.ID,
		fmt.Sprintf("```⭕️ タスク追加: %s 優先度： %d (%s)```", title, priorityID, priorityEmoji[priorityID]),
	); err != nil {
		return errors.NewAppError("taskAddedButReplyFailed", err)
	}

	return nil
}

func HandleList(s *discordgo.Session, m *discordgo.MessageCreate) error {
	tasks, err := service.GetTaskService(m.Author.ID)
	if err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, "```❌ タスク取得失敗```"); rerr != nil {
			return errors.NewAppError("failedToGetTasks", rerr)
		}
		return err
	}

	if len(tasks) == 0 {
		if err := replyToUser(s, m.ChannelID, m.Author.ID, "```📭 タスクが登録されていません```"); err != nil {
			return errors.NewAppError("noTasks", err)
		}
		return nil
	}

	var msg strings.Builder
	msg.WriteString("今日のTodoです！\n```")
	completedFlag := false
	for i, task := range tasks {
		if task.Status == "pending" {
			if i == 0 {
				msg.WriteString("📝 未完了のタスク\n")
			}
			msg.WriteString(fmt.Sprintf("%s ⌛️ [%02d] %s\n", priorityEmoji[task.PriorityID], i, task.Title))
		} else if task.Status == "completed" {
			if !completedFlag {
				msg.WriteString("\n✅ 完了済みのタスク\n")
				completedFlag = true
			}
			msg.WriteString(fmt.Sprintf("✅ [%02d] %s\n", i, task.Title))
		}
	}
	msg.WriteString("```")

	if err := replyToUser(s, m.ChannelID, m.Author.ID, msg.String()); err != nil {
		return errors.NewAppError("afterFetchedList", err)
	}
	return nil
}

func HandleComplete(s *discordgo.Session, m *discordgo.MessageCreate, content string) error {
	arg := strings.TrimPrefix(content, fmt.Sprintf("%sdone ", COMMAND_PREFIX))
	DoneTaskNumber, err := strconv.Atoi(arg)
	if err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, "```❌ 数字を指定してください```"); rerr != nil {
			return errors.NewAppError("invalidNumber", rerr)
		}
		return errors.NewAppError("invalidNumber", err)
	}

	if err := service.CompleteTaskService(m.Author.ID, DoneTaskNumber); err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ %s```", err.Error())); rerr != nil {
			return errors.NewAppError("HandleComplete.completeFailed.replyToUser", rerr)
		}
		return err
	}

	tasks, err := service.GetTaskService(m.Author.ID)
	if err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, "```✅ タスク完了！\n⚠️ 残りのタスク取得に失敗しました```"); rerr != nil {
			return errors.NewAppError("afterComplete", rerr)
		}
		return errors.NewAppError("afterComplete", err)
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

	if err := replyToUser(s, m.ChannelID, m.Author.ID, msg.String()); err != nil {
		return errors.NewAppError("afterMarkedCompleted", err)
	}
	return nil
}

func HandleDelete(s *discordgo.Session, m *discordgo.MessageCreate, content string) error {
	arg := strings.TrimPrefix(content, fmt.Sprintf("%sdelete ", COMMAND_PREFIX))
	DeleteNumber, err := strconv.Atoi(arg)
	if err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, "```❌ 数字を指定してください```"); rerr != nil {
			return errors.NewAppError("invalidNumber", rerr)
		}
		return errors.NewAppError("invalidNumber", err)
	}

	if err := service.DeleteTaskService(m.Author.ID, DeleteNumber); err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ %s```", err.Error())); rerr != nil {
			return errors.NewAppError("deleteFailed", rerr)
		}
		return err
	}

	if err := replyToUser(s, m.ChannelID, m.Author.ID, "```⭕️ タスク削除しました```"); err != nil {
		return errors.NewAppError("afterRemoved", err)
	}
	return nil
}

func HandleChat(s *discordgo.Session, m *discordgo.MessageCreate, content string) error {
	arg := strings.TrimPrefix(content, fmt.Sprintf("%schat ", COMMAND_PREFIX))
	if len(strings.TrimSpace(arg)) == 0 {
		if err := replyToUser(s, m.ChannelID, m.Author.ID, "```❌ メッセージを入力してください```"); err != nil {
			return errors.NewAppError("emptyMessage", err)
		}
		return errors.NewAppError("emptyMessage", fmt.Errorf("empty chat message"))
	}

	if err := s.ChannelTyping(m.ChannelID); err != nil {
		return errors.NewAppError("ChannelTyping", err)
	}

	reply, err := service.ChatWithContext(m.Author.ID, arg)
	if err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ %s```", err.Error())); rerr != nil {
			return errors.NewAppError("chatFailed", rerr)
		}
		return err
	}

	if err := replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```\n%s\n```", reply)); err != nil {
		return errors.NewAppError("replyToUser", err)
	}
	return nil
}

func HandleReset(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if strings.HasPrefix(m.Content, fmt.Sprintf("%sreset all", COMMAND_PREFIX)) {
		resetAllConfirm[m.Author.ID] = time.Now().Add(10 * time.Minute)
		if err := replyToUser(
			s, m.ChannelID, m.Author.ID,
			"```⚠️ 本当に全タスク（過去含む）を削除しますか？\n削除するには '!confirm reset' と入力してください。（10分以内）```",
		); err != nil {
			return errors.NewAppError("resetAllConfirm", err)
		}
		return nil
	}

	count, err := service.ResetTodayTasks(m.Author.ID)
	if err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ 今日のリセット失敗: %s```", err.Error())); rerr != nil {
			return errors.NewAppError("resetToday", rerr)
		}
		return err
	}

	if err := replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```✅ 今日のタスクを %d 件削除しました```", count)); err != nil {
		return errors.NewAppError("afterReset", err)
	}
	return nil
}

func HandleConfirm(s *discordgo.Session, m *discordgo.MessageCreate) error {
	userID := m.Author.ID
	expiry, ok := resetAllConfirm[userID]
	if !ok || time.Now().After(expiry) {
		delete(resetAllConfirm, userID)
		if err := replyToUser(s, m.ChannelID, userID, "```⚠️ '!reset all' の確認時間が切れました。再度実行してください。```"); err != nil {
			return errors.NewAppError("expired", err)
		}
		return errors.NewAppError("expired", fmt.Errorf("confirmation expired"))
	}

	count, err := service.ResetAllTasks(userID)
	if err != nil {
		if rerr := replyToUser(s, m.ChannelID, userID, fmt.Sprintf("```❌ 全削除に失敗しました: %s```", err.Error())); rerr != nil {
			return errors.NewAppError("resetAll", rerr)
		}
		return err
	}

	delete(resetAllConfirm, userID)
	if err := replyToUser(s, m.ChannelID, userID, fmt.Sprintf("```✅ 全タスクを %d 件削除しました```", count)); err != nil {
		return errors.NewAppError("afterDeleted", err)
	}
	return nil
}

func HandleEdit(s *discordgo.Session, m *discordgo.MessageCreate, content string) error {
	arg := strings.TrimPrefix(content, fmt.Sprintf("%sedit ", COMMAND_PREFIX))
	fields := strings.Fields(arg)
	if len(fields) < 2 {
		if err := replyToUser(s, m.ChannelID, m.Author.ID,
			"```⚠️ コマンドの形式が正しくありません。\n例: `!edit 1 <title/優先度>` or `!edit 1 title 優先度` ```"); err != nil {
			return errors.NewAppError("invalidFormat.replyToUser", err)
		}
		return errors.NewAppError("invalidFormat", fmt.Errorf("invalid command format"))
	}

	IndexNumber, err := strconv.Atoi(fields[0])
	if err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, "```❌ 数字を指定してください```"); rerr != nil {
			return errors.NewAppError("invalidNumber", rerr)
		}
		return errors.NewAppError("invalidNumber", err)
	}

	// validate input
	params := fields[1:]
	var newPriority *int
	var newTitle string

	if pid, ok := priorityMap[params[len(params)-1]]; ok {
		// paramの末尾が優先度指定なら設定
		newPriority = &pid
	}

	titleEnd := len(params)
	if newPriority != nil {
		// 優先度が指定されている場合はタイトルの終端を調整
		titleEnd -= 1
	}
	if titleEnd > 0 {
		// タイトルが存在する場合のみ設定
		newTitle = strings.Join(params[0:titleEnd], " ")
	}

	if err := service.UpdateTaskService(m.Author.ID, IndexNumber, newTitle, newPriority); err != nil {
		if rerr := replyToUser(s, m.ChannelID, m.Author.ID, fmt.Sprintf("```❌ タスクの編集に失敗しました: %s```", err.Error())); rerr != nil {
			return errors.NewAppError("updateFailed", rerr)
		}
		return err
	}

	if err := replyToUser(s, m.ChannelID, m.Author.ID, "```✅ 指定されたToDoを編集しました```"); err != nil {
		return errors.NewAppError("afterUpdated", err)
	}
	return nil
}

func HandleHelp(s *discordgo.Session, m *discordgo.MessageCreate) error {
	helpText := "**📋 Self-Management Bot コマンド一覧**\n" +
		"以下のコマンドを使って、タスクの管理やAIとの対話ができます！\n\n" +
		"```" +
		"✅ タスク管理\n" +
		"!add <タスク名> [P1~P4]        : タスクを追加（例: !add 宿題 P1）\n" +
		"!list                         : 今日のタスクを一覧表示\n" +
		"!done <番号>                  : 指定タスクを完了扱いに\n" +
		"!edit <番号> <内容> [P1~P4]   : 内容や優先度を編集\n" +
		"!delete <番号>                : 指定タスクを削除\n\n" +
		"♻️ タスク全削除（慎重に）\n" +
		"!reset                        : 今日のタスクを全削除\n" +
		"!reset all                    : 全タスクを削除（確認付き）\n" +
		"!confirm reset                : 全削除を確定\n\n" +
		"🤖 AI機能\n" +
		"!chat <メッセージ>            : AIと会話（モチベ維持や相談）\n\n" +
		"❓ ヘルプ\n" +
		"!help                         : このヘルプを再表示\n" +
		"```"
	if err := replyToUser(s, m.ChannelID, m.Author.ID, helpText); err != nil {
		return errors.NewAppError("HandleHelp.replyToUser", err)
	}
	return nil
}
