package wh

import (
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var (
	leadCh chan types.JID = make(chan types.JID)
	msgCh  chan Message   = make(chan Message)
)

type Message struct {
	Text       string
	Path       string
	Type       string
	UserId     int
	LeadPhone  string
	IsFromUser bool
}

func ListenLead() string {
	jid := <-leadCh
	return jid.User
}

func ListenMsg() Message {
	return <-msgCh
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

func HandleChatEvents(userId int) func(any) {
	return func(evt any) {
		msg, ok := evt.(*events.Message)
		if !ok {
			return
		}
		if msg.Info.IsGroup {
			return
		}
		msgCh <- Message{
			UserId:     userId,
			Text:       *msg.Message.Conversation,
			IsFromUser: msg.Info.IsFromMe,
			LeadPhone:  msg.Info.Chat.User,
			Path:       "",
			Type:       "text",
		}
	}
}
