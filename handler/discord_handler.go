package handler

import (
    "fmt"
    "strings"

    "github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }

    content := strings.TrimSpace(m.Content)

    if strings.HasPrefix(content, "!add ") {
        taskTitle := strings.TrimPrefix(content, "!add ")
        fmt.Println("タスク登録:", taskTitle)
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ タスク '%s' を登録しました！", taskTitle))
    }
}
