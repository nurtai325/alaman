package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

func AddSession(user User) *http.Cookie {
	b := make([]byte, 32)
	rand.Read(b)
	sessionId := base64.URLEncoding.EncodeToString(b)
	sessions.mu.Lock()
	defer sessions.mu.Unlock()
	user.Valid = true
	sessions.s[sessionId] = sessionInfo{
		id:      sessionId,
		user:    user,
		expires: time.Now().AddDate(0, 0, sessionExpDays),
	}
	return newCookie(sessionId)
}

func DeleteSession(r *http.Request) *http.Cookie {
	sessionInfo, _ := getSession(r)
	sessions.mu.Lock()
	defer sessions.mu.Unlock()
	delete(sessions.s, sessionInfo.id)
	return newCookie("")
}

func newCookie(sessionId string) *http.Cookie {
	return &http.Cookie{
		Name:   sessionCookieName,
		Value:  sessionId,
		MaxAge: sessionCookieMaxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
