package wh

import (
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var (
	leadCh chan types.JID = make(chan types.JID)
)

func ListenLead() string {
	jid := <-leadCh
	return jid.User
}

func HandleLeadEvents(evt interface{}) {
	e, ok := evt.(*events.Message)
	if !ok {
		return
	}
	if e.Info.IsFromMe || e.Info.IsGroup {
		return
	}
	leadCh <- e.Info.Sender
}
