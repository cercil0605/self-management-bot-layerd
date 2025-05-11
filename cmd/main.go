package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"self-management-bot/client"
	"self-management-bot/db"
	"self-management-bot/handler"
	"time"
)

func main() {
	// env read
	err := godotenv.Load(".env") // 相対パスに注意
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}
	token := os.Getenv("DISCORD_BOT_TOKEN")
	// Docker PostgreSQL 起動
	if client.IsDockerPostgresRunning() {
		log.Println("✅ PostgreSQL は既に起動済み")
	} else {
		log.Println("🚀 PostgreSQL 起動開始")
		if err := client.RunDockerSQL(); err != nil {
			log.Fatal("❌ PostgreSQL 起動失敗:", err)
		}
		time.Sleep(2 * time.Second)
	}
	// Connect DB
	if err := db.Init(); err != nil {
		log.Fatal("❌ DB 初期化失敗:", err)
	}
	log.Println("✅ DB 初期化成功")
	// session with Discord
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("❌ Error creating Discord session,", err)
	}
	dg.AddHandler(handler.MessageCreate)
	log.Println("✅ Discordセッション成功")
	// connect with Discord
	err = dg.Open()
	if err != nil {
		log.Fatal("❌ Error opening Discord connection,", err)
	}
	log.Println("✅ Discord接続成功")

	defer dg.Close()
	// パッチ処理
	handler.StartResetConfirmCleaner()
	handler.StartFixedReminderSender(dg)
	// boot LLM
	err = client.StartLLM()
	if err != nil {
		log.Fatal("❌ Error opening LLM connection,", err)
	}

	defer client.StopLLM()
	log.Println("✅ Bot is now running. Press CTRL+C to exit.")
	select {}
}
