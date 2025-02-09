package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/db/repository"
	"github.com/nurtai325/alaman/internal/wh"
)

var (
	ErrAlreadyPaired = errors.New("user is already paired")
)

type LeadWh struct {
	Id        int
	Name      string
	Phone     string
	Jid       string
	CreatedAt time.Time
}

func getSLeadWh(l repository.LeadWh) LeadWh {
	return LeadWh{
		Id:        int(l.ID),
		Name:      l.Name,
		Phone:     l.Phone,
		Jid:       l.Jid.String,
		CreatedAt: l.CreatedAt.Time,
	}
}

func (s *Service) GetLeadWhs(ctx context.Context, offset, limit int) ([]LeadWh, error) {
	if offset < 0 {
		return nil, ErrInvalidOffset
	} else if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	leadWhs, err := s.queries.GetLeadWhs(ctx, repository.GetLeadWhsParams{
		Offset: int64(offset),
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, err
	}
	sLeadWhs := make([]LeadWh, 0, len(leadWhs))
	for _, leadWh := range leadWhs {
		sLeadWhs = append(sLeadWhs, getSLeadWh(leadWh))
	}
	return sLeadWhs, nil
}

func (s *Service) GetLeadWh(ctx context.Context, id int) (LeadWh, error) {
	p, err := s.queries.GetLeadWh(ctx, int32(id))
	if err != nil {
		return LeadWh{}, err
	}
	return getSLeadWh(p), nil
}

func (s *Service) InsertLeadWh(ctx context.Context, name, phone string) (LeadWh, error) {
	if !validPhone(phone) {
		return LeadWh{}, ErrInvalidPhone
	}
	p, err := s.queries.InsertLeadWh(ctx, repository.InsertLeadWhParams{
		Name:  name,
		Phone: phone,
	})
	if err != nil {
		return LeadWh{}, errors.Join(ErrInternal, err)
	}
	return getSLeadWh(p), nil
}

func (s *Service) UpdateLeadWh(ctx context.Context, id int, name, phone string) (LeadWh, error) {
	if !validPhone(phone) {
		return LeadWh{}, ErrInvalidPhone
	}
	p, err := s.queries.UpdateLeadWh(ctx, repository.UpdateLeadWhParams{
		ID:    int32(id),
		Name:  name,
		Phone: phone,
	})
	if err != nil {
		return LeadWh{}, errors.Join(ErrInternal, err)
	}
	return getSLeadWh(p), nil
}

func (s *Service) DeleteLeadWh(ctx context.Context, id int) (LeadWh, error) {
	if err := validId(id); err != nil {
		return LeadWh{}, err
	}
	p, err := s.queries.DeleteLeadWh(ctx, int32(id))
	if err != nil {
		return LeadWh{}, err
	}
	return getSLeadWh(p), nil
}

func (s *Service) ConnectLeadWh(ctx context.Context, id int, jid string) (LeadWh, error) {
	if err := validId(id); err != nil {
		return LeadWh{}, err
	}
	p, err := s.queries.ConnectLeadWh(ctx, repository.ConnectLeadWhParams{
		ID: int32(id),
		Jid: pgtype.Text{
			String: jid,
			Valid:  true,
		},
	})
	if err != nil {
		return LeadWh{}, err
	}
	return getSLeadWh(p), nil
}

func (s *Service) GetLeadWhQr(phone string) (string, error) {
	if !validPhone(phone) {
		return "", ErrInvalidPhone
	}
	imagePath, err := wh.StartPairing(phone[1:], wh.LeadEventsHandler)
	if err != nil && errors.Is(err, wh.ErrAlreadyPaired) {
		return "", ErrAlreadyPaired
	}
	return imagePath, err
}

func (s *Service) ConnectAllWh() error {
	leadWhs, err := s.GetLeadWhs(context.Background(), 0, 1000)
	if err != nil {
		return err
	}
	for _, leadWh := range leadWhs {
		_, err := wh.Connect(leadWh.Jid, wh.LeadEventsHandler)
		if err != nil {
			log.Println(err)
		}
	}
	users, err := s.GetUsers(context.Background(), 0, 1000)
	if err != nil {
		return err
	}
	for i, user := range users {
		client, err := wh.Connect(user.Jid, wh.ChatEventsHandler(user.Id))
		if err != nil {
			log.Println(err)
		}
		if i == 0 {
			wh.SetDefaultClient(client)
		}
	}
	return nil
}
