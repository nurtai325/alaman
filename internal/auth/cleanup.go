package auth

import "time"

func Cleanup() {
	for {
		time.Sleep(time.Hour * 12)
		sessions.mu.Lock()
		now := time.Now()
		for sessionId, info := range sessions.s {
			if info.expires.Before(now) {
				delete(sessions.s, sessionId)
			}
		}
		sessions.mu.Unlock()
	}
}
