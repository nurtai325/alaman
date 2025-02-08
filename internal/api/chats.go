package api

import (
	"net/http"
	"strconv"

	"github.com/nurtai325/alaman/internal/service"
)

type chatsData struct {
	Chats []service.Chat
}

func (app *app) handleChatsGet(w http.ResponseWriter, r *http.Request) {
	chats, err := app.service.GetChats(r.Context(), 0, pagesLimit)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tChats, "/pages/chats", layoutData{
		BarsData: barsData{
			Page:     "chats",
			PageName: "Чаттар",
			Pages:    getPage(r),
		},
		User: app.service.GetAuthUser(r),
		Data: chatsData{
			Chats: chats,
		},
	})
}

type messagesData struct {
	Phone    string
	Sender   string
	Messages []service.Message
}

func (app *app) handleMessagesGet(w http.ResponseWriter, r *http.Request) {
	chatIdStr := r.PathValue("id")
	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		app.error(w, err)
		return
	}
	messages, err := app.service.GetMessages(r.Context(), chatId)
	if err != nil {
		app.error(w, err)
		return
	}
	chat, err := app.service.GetChat(r.Context(), chatId)
	if err != nil {
		app.error(w, err)
		return
	}
	app.execute(w, tMessages, "", messagesData{
		Phone:    chat.LeadPhone,
		Sender:   chat.UserName,
		Messages: messages,
	})
}
