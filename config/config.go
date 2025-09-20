package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config はアプリケーション全体の設定を保持します。
type Config struct {
	DiscordToken string
	GeminiApiKey string
}

// Cfg はロードされた設定を保持するグローバル変数です。
var Cfg *Config

// LoadConfig は環境変数または.envファイルから設定を読み込みます。
func LoadConfig() {
	// .envファイルから環境変数を読み込む（ファイルが存在しない場合もエラーにしない）
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("環境変数 'DISCORD_TOKEN' が設定されていません。")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("環境変数 'GEMINI_API_KEY' が設定されていません。")
	}

	Cfg = &Config{
		DiscordToken: token,
		GeminiApiKey: apiKey,
	}
}
