package wh

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

var (
	clients   map[string]*cliWh = make(map[string]*cliWh)
	container *sqlstore.Container
)

type cliWh struct {
	mu sync.Mutex
	c  *whatsmeow.Client
}

func InitContainer(dbConn *sql.DB) error {
	container = sqlstore.NewWithDB(dbConn, "postgres", nil)
	err := container.Upgrade()
	if err != nil {
		return fmt.Errorf("error making new sql whatsapp container: %w", err)
	}
	return nil
}

func SendMessage(ctx context.Context, from, to, text string) error {
	return nil
}

func Connect(jidStr string, eventHandler func(any)) error {
	jid, err := types.ParseJID(jidStr)
	if err != nil {
		return err
	}
	client, ok := clients[jid.User]
	if ok {
		client.c.AddEventHandler(eventHandler)
		return nil
	}
	device, err := container.GetDevice(jid)
	if err != nil {
		return err
	}
	if device == nil {
		return fmt.Errorf("jid %s: %w", jid, ErrDeviceNotFound)
	}
	newClient := whatsmeow.NewClient(device, nil)
	err = newClient.Connect()
	if err != nil {
		return err
	}
	newClient.AddEventHandler(eventHandler)
	clients[jid.User] = &cliWh{c: newClient}
	return nil
}

func GetJid(phone string) string {
	return clients[phone].c.Store.ID.String()
}
