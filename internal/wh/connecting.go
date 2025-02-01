package wh

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
)

type pairingData struct {
	ImagePath string
}

func GetQr() (pairingData, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.c != nil {
		return pairingData{}, nil
	}
	device, err := container.GetFirstDevice()
	if err != nil || device == nil {
		return pairingData{}, err
	}
	newClient := whatsmeow.NewClient(device, nil)
	newClient.AddEventHandler(handleEvents)
	client.c = newClient
	qrCh, _ := client.c.GetQRChannel(context.Background())
	err = client.c.Connect()
	if err != nil {
		return pairingData{}, fmt.Errorf("error connecting to whatsapp websocket: %w", err)
	}
	for evt := range qrCh {
		if evt.Event != "code" {
			continue
		}
		imagePath := ""
		go func(qrCh <-chan whatsmeow.QRChannelItem) {
			for evt := range qrCh {
				if evt.Error != nil {
					err := fmt.Errorf("image path: %s event: %s: %w", imagePath, evt.Event, evt.Error)
					log.Println(err)
					return
				}
				if evt.Event == "timeout" {
					err := fmt.Errorf("qr channel timed out. image path: %s event: %s: %w", imagePath, evt.Event, evt.Error)
					log.Println(err)
					return
				}
				if evt.Event == "success" {
					if err != nil {
						log.Println(err)
					}
					return
				}
			}
			log.Println(fmt.Errorf("qr channel closed unexpectedly"))
		}(qrCh)
		imagePath = fmt.Sprintf("/assets/qr/%d.qr.png", time.Now().UnixNano())
		err = qrcode.WriteFile(evt.Code, qrcode.Medium, 512, "."+imagePath)
		if err != nil {
			return pairingData{}, fmt.Errorf("error generating qr code image: %w", err)
		}
		return pairingData{ImagePath: imagePath}, nil
	}
	return pairingData{}, fmt.Errorf("unsuccesfull pairing")
}
