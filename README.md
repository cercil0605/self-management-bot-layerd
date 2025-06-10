# Self-Management Bot

🎯 **日々のタスク管理を Discord 上で行える自己管理Bot**  
⌛️ タスクの追加・完了・削除・一覧機能に加え、AIがあなたのメンタルコーチになります。  
🧠 チャットを通じてモチベーション維持や振り返りをサポートします。

`!chat` コマンドでは LLM との対話が可能です。ローカルの ELYZA モデルまたは Google Gemini API のどちらかを利用できます。

---

## 🎮 対応プラットフォーム

- Discord

---

## 🚀 コマンド一覧

| コマンド                            | 説明                 |
|---------------------------------|--------------------|
| `!add <内容> <優先度>`               | タスクを追加，4段階の優先度設定可能 |
| `!list`                         | 当日タスクを一覧表示         |
| `!edit <番号> <タイトル> <優先度>`       | タスクのタイトルを編集        |
| `!done <番号>`                    | 指定した番号のタスクを完了      |
| `!delete <番号>`                  | 指定した番号のタスクを削除      |
| `!chat <内容>`                    | LLMとの会話                 |
| `!reset`                        | 当日分のタスクを全削除        |
| `!reset all` → `!confirm reset` | 全タスクを完全に削除         |

---

## 🛠️ 技術スタック

- **Language**: Go 1.20+
- **Discord API**: [`discordgo`](https://github.com/bwmarrin/discordgo)
- **Database**: PostgreSQL + [`sqlx`](https://github.com/jmoiron/sqlx)
- **LLM API**:  [`Ollama`](https://github.com/ollama/ollama) / [Google Gemini](https://ai.google.dev/)
- **Infra**: Docker（PostgreSQL用）

### 環境変数

`LLM_PROVIDER` に `elyza` (デフォルト) または `gemini` を指定します。Gemini を使用する場合は `GEMINI_API_KEY` も設定してください。
