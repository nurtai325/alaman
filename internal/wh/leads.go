package wh

import (
	"log"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var leadCh chan types.JID = make(chan types.JID)

func ListenLead() types.JID {
	jid := <-leadCh
	return jid
}

func handleEvents(evt interface{}) {
	e, ok := evt.(*events.Message)
	if !ok {
		return
	}
	if e.Info.IsFromMe || e.Info.IsGroup {
		return
	}
	client.mu.Lock()
	defer client.mu.Unlock()
	c, err := client.c.Store.Contacts.GetContact(e.Info.Sender)
	if err != nil {
		log.Println(err)
		return
	}
	if c.FullName != "" {
		return
	}
	jid := e.Info.Sender
	leadCh <- jid
}
