package wh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
)

var (
	ErrAlreadyPaired = errors.New("user is already paired")
	ErrDeviceNotFound = errors.New("wh with this phone is not present in store")
)

func StartPairing(phone string, eventHandler func(evt any)) (string, error) {
	_, found := clients[phone]
	if found {
		return "", ErrAlreadyPaired
	}
	device := container.NewDevice()
	newClient := whatsmeow.NewClient(device, nil)
	newClient.AddEventHandler(eventHandler)
	qrCh, err := newClient.GetQRChannel(context.Background())
	if err != nil {
		return "", fmt.Errorf("error getting qr channel: %w", err)
	}
	err = newClient.Connect()
	if err != nil {
		return "", fmt.Errorf("error connecting to whatsapp websocket: %w", err)
	}
	for evt := range qrCh {
		if evt.Event != "code" {
			continue
		}
		imagePath := ""
		go waitPairing(qrCh, newClient, phone)
		imagePath = fmt.Sprintf("/assets/qr/%d.qr.png", time.Now().UnixNano())
		err = qrcode.WriteFile(evt.Code, qrcode.Medium, 512, "."+imagePath)
		if err != nil {
			return "", fmt.Errorf("error generating qr code image: %w", err)
		}
		return imagePath, nil
	}
	return "", fmt.Errorf("unsuccesfull pairing")
}

func waitPairing(qrCh <-chan whatsmeow.QRChannelItem, client *whatsmeow.Client, phone string) {
	for evt := range qrCh {
		if evt.Error != nil {
			err := fmt.Errorf("event: %s: %w", evt.Event, evt.Error)
			log.Println(err)
			return
		}
		if evt.Event == "timeout" {
			err := fmt.Errorf("qr channel timed out. event: %s: %w", evt.Event, evt.Error)
			log.Println(err)
			return
		}
		if evt.Event == "success" {
			clients[phone] = &cliWh{c: client}
			return
		}
	}
	log.Println(fmt.Errorf("qr channel closed unexpectedly"))
}
