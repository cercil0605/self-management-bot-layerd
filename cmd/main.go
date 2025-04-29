package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"self-management-bot/db"
	"self-management-bot/handler"
)

func main() {
	// env read
	err := godotenv.Load(".env") // 相対パスに注意
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}
	token := os.Getenv("DISCORD_BOT_TOKEN")
	// DB Connection
	if err := db.Init(); err != nil {
		panic(err)
	}
	// session with discord
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("❌ Error creating Discord session,", err)
	}

	dg.AddHandler(handler.MessageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal("❌ Error opening connection,", err)
	}
	defer dg.Close()

	log.Println("✅ Bot is now running. Press CTRL+C to exit.")
	select {}
}
