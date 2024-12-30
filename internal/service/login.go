package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/nurtai325/alaman/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidLoginInfo = errors.New("Телефон номері немесе құпиясөз қате")
)

func (s *Service) Login(ctx context.Context, phone, password string) (*http.Cookie, error) {
	if !validPassword(password) || !validPhone(phone) {
		return nil, ErrInvalidLoginInfo
	}
	user, err := s.queries.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = ErrInvalidLoginInfo
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidLoginInfo
	}
	sessionCookie := auth.AddSession(auth.User{
		Id:    int(user.ID),
		Phone: user.Phone,
		Name:  user.Name,
		Valid: true,
	})
	return sessionCookie, nil
}

func validPhone(phone string) bool {
	if phone == "" {
		return false
	} else if len(phone) != 12 {
		return false
	} else if rune(phone[0]) != '+' {
		return false
	}
	phone = phone[1:]
	for _, r := range phone {
		if r <= 47 || r >= 58 {
			return false
		}
	}
	return true
}

func validPassword(password string) bool {
	if password == "" || len(password) > 72 {
		return false
	} else if len(password) < 8 {
		return false
	}
	return true
}
