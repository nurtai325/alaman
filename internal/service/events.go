package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/db/repository"
	"github.com/nurtai325/alaman/internal/wh"
)

func ListenNewLeads(s *Service) {
	for {
		phone := wh.ListenLead()
		phone = "+" + phone
		lead, err := s.GetLeadByPhone(context.Background(), phone)
		if err != nil && !errors.Is(err, ErrNotFound) {
			log.Println(err)
			continue
		}
		if lead.Id != 0 {
			continue
		}
		_, err = s.InsertLead(context.Background(), phone)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func ListenNewMessages(s *Service) {
	for {
		msg := wh.ListenMsg()
		phone := "+" + msg.LeadPhone
		lead, err := s.queries.GetLeadByPhone(context.Background(), phone)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			log.Println(err)
			continue
		}
		chat, err := s.queries.GetChatByLeadId(context.Background(), lead.ID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				log.Println(err)
				continue
			}
			chat, err = s.queries.InsertChat(context.Background(), repository.InsertChatParams{
				LeadID: lead.ID,
				UserID: int32(msg.UserId),
			})
			if err != nil {
				log.Println(err)
				continue
			}
		}
		_, err = s.queries.InsertMessage(context.Background(), repository.InsertMessageParams{
			Text: pgtype.Text{
				Valid:  true,
				String: msg.Text,
			},
			Path: pgtype.Text{
				Valid:  true,
				String: msg.Path,
			},
			Type:        msg.Type,
			IsSent:      msg.IsFromUser,
			AudioLength: int32(msg.AudioLength),
			ChatID:      chat.ID,
		})
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = s.queries.UpdateChat(context.Background(), repository.UpdateChatParams{
			ID: chat.ID,
			UpdatedAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
