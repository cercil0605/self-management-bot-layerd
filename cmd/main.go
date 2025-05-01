package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"self-management-bot/client"
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
	// boot Docker
	err = client.RunDockerSQL()
	if err != nil {
		log.Fatal("❌ Error opening Docker connection,", err)
	}
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
		log.Fatal("❌ Error opening Discord connection,", err)
	}
	defer dg.Close()

	// boot LLM
	err = client.StartLLM()
	if err != nil {
		log.Fatal("❌ Error opening LLM connection,", err)
	}
	defer client.StopLLM()

	log.Println("✅ Bot is now running. Press CTRL+C to exit.")
	select {}
}
