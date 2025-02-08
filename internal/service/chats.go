package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/db/repository"
)

type Chat struct {
	Id        int
	Messages  []Message
	LeadId    int
	UserId    int
	CreatedAt time.Time
}

type msgType string

const (
	textMsg  msgType = "text"
	audioMsg msgType = "audio"
	photoMsg msgType = "photo"
	videoMsg msgType = "video"
)

type Message struct {
	Id        int
	Text      string
	Path      string
	Type      msgType
	IsSent    bool
	ChatId    int
	CreatedAt time.Time
}

func (s *Service) GetChats(ctx context.Context, offset, limit int) ([]Chat, error) {
	chats, err := s.queries.GetChats(ctx, repository.GetChatsParams{
		Offset: int64(offset),
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, err
	}
	sChats := make([]Chat, 0, len(chats))
	for _, chat := range chats {
		messages, err := s.GetMessages(ctx, int(chat.ID))
		if err != nil {
			return nil, err
		}
		sChats = append(sChats, getSChat(chat, messages))
	}
	return sChats, nil
}

func (s *Service) GetMessages(ctx context.Context, chatId int) ([]Message, error) {
	messages, err := s.queries.GetMessages(ctx, int32(chatId))
	if err != nil {
		return nil, err
	}
	sMessages := make([]Message, 0, len(messages))
	for _, msg := range messages {
		sMessages = append(sMessages, getSMessage(msg))
	}
	return sMessages, nil
}

func getSChat(chat repository.Chat, messages []Message) Chat {
	return Chat{
		Id:        int(chat.ID),
		Messages:  messages,
		LeadId:    int(chat.LeadID),
		UserId:    int(chat.UserID),
		CreatedAt: chat.CreatedAt.Time,
	}
}

func (s *Service) InsertChat(ctx context.Context, leadId, userId int) (Chat, error) {
	chat, err := s.queries.InsertChat(ctx, repository.InsertChatParams{
		LeadID: int32(leadId),
		UserID: int32(userId),
	})
	if err != nil {
		return Chat{}, err
	}
	return getSChat(chat, nil), err
}

func (s *Service) InsertMessage(ctx context.Context, text, path string, msgtype msgType, isSent bool, audioLength, chatId int) (Message, error) {
	msg, err := s.queries.InsertMessage(ctx, repository.InsertMessageParams{
		Text: pgtype.Text{
			Valid:  true,
			String: text,
		},
		Path: pgtype.Text{
			Valid:  true,
			String: path,
		},
		Type:        string(msgtype),
		IsSent:      isSent,
		ChatID:      int32(chatId),
		AudioLength: int32(audioLength),
	})
	if err != nil {
		return Message{}, err
	}
	return getSMessage(msg), nil
}

func getSMessage(msg repository.Message) Message {
	return Message{
		Id:        int(msg.ID),
		Text:      msg.Text.String,
		Path:      msg.Path.String,
		Type:      msgType(msg.Type),
		IsSent:    msg.IsSent,
		ChatId:    int(msg.ChatID),
		CreatedAt: msg.CreatedAt.Time,
	}
}
