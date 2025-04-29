package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"self-management-bot/service"
	"strings"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := strings.TrimSpace(m.Content)

	if strings.HasPrefix(content, "!add ") {
		title := strings.TrimPrefix(content, "!add ")
		err := service.AddTaskService(m.Author.ID, title)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\n "+"❌ タスク登録失敗")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@"+m.Author.ID+">\n "+"✅ タスク追加: %s", title))
	}
}
