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
		jid := wh.ListenLead()
		lead, err := s.GetLeadByPhone(context.Background(), "+"+jid.User)
		if err != nil && !errors.Is(err, ErrNotFound) {
			log.Println(err)
			continue
		}
		if lead.Id != 0 || lead.CreatedAt.After(time.Now().AddDate(0, 0, -7)) {
			continue
		}
		_, err = s.InsertLead(context.Background(), "+"+jid.User)
		if err != nil {
			log.Println(err)
			continue
		}
		err = wh.Archive(context.Background(), jid)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
