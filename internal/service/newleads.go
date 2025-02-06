package service

import (
	"context"
	"errors"
	"log"
	"time"

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
		if lead.Id != 0 || lead.CreatedAt.After(time.Now().AddDate(0, 0, -7)) {
			continue
		}
		_, err = s.InsertLead(context.Background(), phone)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
