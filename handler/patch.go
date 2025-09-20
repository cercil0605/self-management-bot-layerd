package handler

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"self-management-bot/service"
	"time"
)

// StartResetConfirmCleaner は、リセット確認の有効期限が切れたものを掃除します。
// 注意: resetAllConfirmへのアクセスはスレッドセーフである必要があります。
func StartResetConfirmCleaner() {
	go func() {
		// 1分ごとに期限切れの確認を削除する
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			// TODO: resetAllConfirmへのアクセスをMutexで保護する
			for userID, expiry := range resetAllConfirm {
				if now.After(expiry) {
					delete(resetAllConfirm, userID)
				}
			}
		}
	}()
}

// StartFixedReminderSender は、指定された時刻にリマインダーを送信します。
// 送信時刻: 6:00, 12:00, 19:00
func StartFixedReminderSender(s *discordgo.Session) {
	go func() {
		// 1分ごとに時刻をチェックするTicker
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for t := range ticker.C {
			// 現在時刻の分が0（正時）の場合のみ処理を検討
			if t.Minute() == 0 {
				hour := t.Hour()
				// 指定された時刻（6時, 12時, 19時）であればリマインダーを送信
				if hour == 6 || hour == 12 || hour == 19 {
					log.Println("✉️ " + "リマインド開始")
					SendReminder(s)
					log.Println("✅ " + "リマインド完了")
				}
			}
		}
	}()
}

// SendReminder は、リマインド対象の全ユーザーにメッセージを送信します。
func SendReminder(s *discordgo.Session) {
	reminders, err := service.FixedTimeReminder()
	if err != nil {
		log.Printf("❌ リマインド取得エラー: %v", err)
		return
	}

	if len(reminders) == 0 {
		log.Println("⚠️ リマインド対象が0件です。送信スキップ")
		return
	}

	for _, reminder := range reminders {
		// DMチャンネルを作成
		channel, err := s.UserChannelCreate(reminder.UserID)
		if err != nil {
			log.Printf("⚠️ DMチャンネル取得失敗 userID=%s: %v", reminder.UserID, err)
			continue // 次のユーザーへ
		}

		// DMを送信
		_, err = s.ChannelMessageSend(channel.ID, reminder.Content)
		if err != nil {
			log.Printf("❌ リマインド送信失敗 userID=%s: %v", reminder.UserID, err)
		}
	}
}
