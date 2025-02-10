package service

import (
	"context"
	"errors"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/auth"
	"github.com/nurtai325/alaman/internal/db/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrMinMax              = errors.New("min id is bigger than max id")
	ErrInvalidPhone        = errors.New("телефон номерінің форматы қате")
	ErrInvalidPassword     = errors.New("құпиясөз форматы қате")
	ErrUserNameTooLong     = errors.New("есім тым үлкен")
	ErrPasswordNotEqual    = errors.New("құпиясөздер тең емес")
	ErrUserWithPhoneExists = errors.New("бұндай номермен қызметкер бар. басқасын таңдаңыз")
)

type User struct {
	Id        int
	Name      string
	Phone     string
	Role      auth.Role
	Jid       string
	CreatedAt time.Time
}

func (user User) GetId() int {
	return user.Id
}

func getSUser(user repository.User) User {
	return User{
		Id:        int(user.ID),
		Name:      user.Name,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Time,
		Role:      auth.Role(user.Role),
		Jid:       user.Jid.String,
	}
}

func (s *Service) GetAuthUser(r *http.Request) User {
	user := auth.GetUser(r)
	return User{
		Id:    user.Id,
		Name:  user.Name,
		Phone: user.Phone,
		Role:  user.Role,
	}
}

func (s *Service) GetUser(ctx context.Context, id int) (User, error) {
	user, err := s.queries.GetUser(ctx, int32(id))
	if err != nil {
		return User{}, err
	}
	return getSUser(user), nil
}

func (s *Service) GetUsers(ctx context.Context, offset, limit int) ([]User, error) {
	if offset < 0 {
		return nil, ErrInvalidOffset
	} else if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	users, err := s.queries.GetUsers(ctx, repository.GetUsersParams{
		Offset: int64(offset),
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, err
	}
	sUsers := make([]User, 0, len(users))
	for _, user := range users {
		sUsers = append(sUsers, getSUser(user))
	}
	return sUsers, nil
}

func (s *Service) GetUsersCount(ctx context.Context) (int, error) {
	count, err := s.queries.GetUsersCount(ctx)
	return int(count), err
}

func (s *Service) UpdateUser(ctx context.Context, id int, name, phone string, role auth.Role) (User, error) {
	if !validPhone(phone) {
		return User{}, ErrInvalidPhone
	} else if err := validUserName(name); err != nil {
		return User{}, err
	} else if err := validId(id); err != nil {
		return User{}, err
	}
	user, err := s.queries.UpdateUser(ctx, repository.UpdateUserParams{
		ID:    int32(id),
		Phone: phone,
		Name:  name,
		Role:  string(role),
	})
	if err != nil {
		return User{}, errors.Join(err, ErrInternal)
	}
	return getSUser(user), nil
}

func (s *Service) InsertUser(ctx context.Context, name, phone, password, checkPassword, jid string, role auth.Role) (User, error) {
	if !validPhone(phone) {
		return User{}, ErrInvalidPhone
	} else if err := validUserName(name); err != nil {
		return User{}, ErrUserNameTooLong
	} else if password != checkPassword {
		return User{}, ErrPasswordNotEqual
	} else if !validPassword(password) {
		return User{}, ErrInvalidPassword
	}
	user, err := s.queries.GetUserByPhone(ctx, phone)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return User{}, err
	}
	if user.ID != 0 {
		return User{}, ErrUserWithPhoneExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.Join(err, ErrInternal)
	}
	user, err = s.queries.InsertUser(ctx, repository.InsertUserParams{
		Name:     name,
		Phone:    phone,
		Password: string(hashedPassword),
		Role:     string(role),
		Jid: pgtype.Text{
			String: jid,
			Valid:  true,
		},
	})
	if err != nil {
		return User{}, errors.Join(err, ErrInternal)
	}
	return getSUser(user), nil
}

func (s *Service) DeleteUser(ctx context.Context, id int) (User, error) {
	if err := validId(id); err != nil {
		return User{}, err
	}
	user, err := s.queries.DeleteUser(ctx, int32(id))
	if err != nil {
		return User{}, err
	}
	return getSUser(user), nil
}

func validUserName(name string) error {
	if utf8.RuneCount([]byte(name)) > 100 {
		return ErrUserNameTooLong
	}
	return nil
}

func validId(id int) error {
	if id <= 0 {
		return errors.Join(ErrInvalidId, ErrInternal)
	}
	return nil
}
