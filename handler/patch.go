package handler

import "time"

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
