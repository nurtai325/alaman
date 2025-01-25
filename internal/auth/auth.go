package auth

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

const (
	sessionCookieName   = "session_id"
	sessionCookieMaxAge = 3600 * 24 * 14
	sessionExpDays      = 7
)

var (
	ErrInvalidRole = errors.New("role doesn't exist")
)

type sessionsMap struct {
	mu sync.Mutex
	s  map[string]sessionInfo
}

type User struct {
	Id    int
	Phone string
	Name  string
	Role  Role
	Valid bool
}

type Role string

const (
	ManagerRole Role = "менеджер"
	LogistRole  Role = "логист"
	AdminRole   Role = "админ"
	RopRole     Role = "роп"
)

func ToRole(role string) (Role, error) {
	switch role {
	case string(ManagerRole):
		return ManagerRole, nil
	case string(LogistRole):
		return LogistRole, nil
	case string(AdminRole):
		return AdminRole, nil
	case string(RopRole):
		return RopRole, nil
	default:
		return "", ErrInvalidRole
	}
}

type sessionInfo struct {
	id      string
	user    User
	expires time.Time
}

var sessions = sessionsMap{
	s: make(map[string]sessionInfo),
}

func IsLogged(r *http.Request) bool {
	_, found := getSession(r)
	return found
}

func GetUser(r *http.Request) User {
	info, _ := getSession(r)
	return info.user
}

func getSession(r *http.Request) (sessionInfo, bool) {
	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return sessionInfo{}, false
	}
	sessions.mu.Lock()
	defer sessions.mu.Unlock()
	info, found := sessions.s[sessionCookie.Value]
	if !found || info.expires.Before(time.Now()) {
		delete(sessions.s, sessionCookie.Value)
		return sessionInfo{}, false
	}
	return info, true
}
