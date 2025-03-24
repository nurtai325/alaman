package wh

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/nurtai325/alaman/internal/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

var (
	clients       map[string]*cliWh = make(map[string]*cliWh)
	defaultClient *whatsmeow.Client
	container     *sqlstore.Container
)

type cliWh struct {
	mu sync.Mutex
	c  *whatsmeow.Client
}

func SetDefaultClient(client *whatsmeow.Client) {
	defaultClient = client
}

func InitContainer(dbConn *sql.DB) error {
	container = sqlstore.NewWithDB(dbConn, "postgres", nil)
	err := container.Upgrade()
	if err != nil {
		return fmt.Errorf("error making new sql whatsapp container: %w", err)
	}
	return nil
}

func SendMessage(ctx context.Context, from, to, text string, isGroup bool) error {
	jid := types.NewJID(to, types.DefaultUserServer)
	if isGroup {
		conf, err := config.New()
		if err != nil {
			return err
		}
		jid, err = types.ParseJID(conf.TABYS_GROUP_ID)
		if err != nil {
			return err
		}
		jid.Server = types.GroupServer
	}
	_, err := defaultClient.SendMessage(ctx, jid, &waE2E.Message{
		Conversation: &text,
	})
	if err != nil {
		fmt.Printf("%+v\n", jid)
		return err
	}
	return nil
}

type whHandler func(*whatsmeow.Client) func(any)

func Connect(jidStr string, eventHandler whHandler) (*whatsmeow.Client, error) {
	jid, err := types.ParseJID(jidStr)
	if err != nil {
		return nil, err
	}
	client, ok := clients[jid.User]
	if ok {
		client.c.AddEventHandler(eventHandler(client.c))
		return client.c, nil
	}
	device, err := container.GetDevice(jid)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, fmt.Errorf("jid %s: %w", jid, ErrDeviceNotFound)
	}
	newClient := whatsmeow.NewClient(device, nil)
	err = newClient.Connect()
	if err != nil {
		return nil, err
	}
	newClient.AddEventHandler(eventHandler(newClient))
	clients[jid.User] = &cliWh{c: newClient}
	return newClient, nil
}

func GetJid(phone string) string {
	return clients[phone].c.Store.ID.String()
}
