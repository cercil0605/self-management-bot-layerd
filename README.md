# Self-Management Bot

🎯 **日々のタスク管理を Discord 上で行える自己管理用Bot**  
⌛️ タスクの追加・完了・削除・一覧機能に加え、AIがあなたのメンタルコーチになります  
⚠️ `!chat`コマンドはLLM APIの実装が必要です．近日デプロイする際にはつけます．

---

## 対応しているアプリ
- Discord
（それ以外は順次対応予定）
---

## 🚀 実行コマンド

- `!add <内容>`：タスクを追加
- `!list`：当日タスクを一覧表示
- `!done <番号>`：指定した番号のタスクを完了
- `!delete <番号>`：指定した番号のタスクを削除
- `!chat <内容>`：過去のタスクと照らし合わせて LLM と会話
- `!reset`：**当日分**のタスクを全削除
- `!reset all` → `!confirm reset`：**全タスク**を削除

---

## 🛠️ 技術スタック

- **Language**: Go 1.20+
- **Discord API**: [`github.com/bwmarrin/discordgo`](https://github.com/bwmarrin/discordgo)
- **DB**: PostgreSQL + [`sqlx`](https://github.com/jmoiron/sqlx)
- **LLM API**:  [`Ollama`](https://github.com/ollama/ollama)
- **RDBMS**: PostgreSQL 
- **Infra**: Docker（DB用）

---

