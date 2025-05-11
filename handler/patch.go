package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"self-management-bot/service"
	"time"
)

// StartResetConfirmCleaner リセット処理のキュー
func StartResetConfirmCleaner() {
	go func() {
		for {
			// 1minごとにresetAllConfirmの期限切れを消す
			time.Sleep(1 * time.Minute)
			now := time.Now()
			for userID, expiry := range resetAllConfirm {
				if now.After(expiry) {
					delete(resetAllConfirm, userID)
				}
			}
		}
	}()
}

// StartFixedReminderSender リマインドを決まった時間に送信する
// TODO 決まった時間に送信できるようにする
// 6:00 12:00 19:00
func StartFixedReminderSender(s *discordgo.Session) {
	go func() {
		for {
			SendReminder(s)
			// テスト用
			time.Sleep(1 * time.Minute)
		}
	}()
}
func SendReminder(s *discordgo.Session) {
	// TODO ユーザー全員分に対応させる
	res, err := service.FixedTimeReminder()
	// とりあえず一人用
	if err != nil {
		fmt.Println("❌ リマインド取得エラー:", err)
		return
	}
	// 要素が空でないかをチェック
	if len(res) == 0 {
		fmt.Println("⚠️ リマインド対象が0件です。送信スキップ")
		return
	}
	var ResForOnePerson = res[0]
	// テスト用 自分のID
	channel, err := s.UserChannelCreate("616909661462855681")
	if err != nil {
		fmt.Printf("⚠️ DMチャンネル取得失敗 userID=%s: %s\n", ResForOnePerson.UserID, err)
		return
	}
	_, err = s.ChannelMessageSend(channel.ID, ResForOnePerson.Content)
	if err != nil {
		fmt.Printf("❌ リマインド送信失敗 userID=%s: %s\n", ResForOnePerson.UserID, err)
	}
}
