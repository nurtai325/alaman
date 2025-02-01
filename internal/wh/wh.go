package wh

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

var (
	client    cliWh
	container *sqlstore.Container
)

type cliWh struct {
	mu sync.Mutex
	c  *whatsmeow.Client
}

func Connect(dbConn *sql.DB) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.c != nil {
		return nil
	}
	newContainer := sqlstore.NewWithDB(dbConn, "postgres", nil)
	err := newContainer.Upgrade()
	if err != nil {
		return fmt.Errorf("error making new sql whatsapp container: %w", err)
	}
	container = newContainer
	device, err := newContainer.GetFirstDevice()
	if err != nil || device == nil {
		return err
	}
	newClient := whatsmeow.NewClient(device, nil)
	err = newClient.Connect()
	if err != nil {
		return err
	}
	newClient.AddEventHandler(handleEvents)
	client.c = newClient
	return nil
}

func Message(ctx context.Context, to, text string) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	_, err := client.c.SendMessage(ctx, types.NewJID(to, types.DefaultUserServer), &waE2E.Message{
		Conversation: proto.String(text),
	})
	if err != nil {
		return fmt.Errorf("whatsapp message sending error to: %s: %w", to, err)
	}
	return nil
}

func GroupMessage(ctx context.Context, text string) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	_, err := client.c.SendMessage(ctx, types.NewJID("120363376293023949", types.GroupServer), &waE2E.Message{
		Conversation: proto.String(text),
	})
	if err != nil {
		return fmt.Errorf("whatsapp group message sending error:  %w", err)
	}
	return nil
}

func Archive(ctx context.Context, jid types.JID) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	return client.c.Store.ChatSettings.PutArchived(jid, true)
}
