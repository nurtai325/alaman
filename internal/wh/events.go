package wh

import (
	"errors"
	"fmt"
	"log"
	"mime"
	"os"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

var (
	ErrUnsupportedMediaType = errors.New("unsupported media type")
)

var (
	leadCh chan string  = make(chan string)
	msgCh  chan Message = make(chan Message)
)

type Message struct {
	Text        string
	Path        string
	Type        string
	UserId      int
	LeadPhone   string
	AudioLength int
	IsFromUser  bool
}

func ListenLead() string {
	return <-leadCh
}

func ListenMsg() Message {
	return <-msgCh
}

func LeadEventsHandler(*whatsmeow.Client) func(any) {
	return func(evt any) {
		e, ok := evt.(*events.Message)
		if !ok {
			return
		}
		if e.Info.IsFromMe || e.Info.IsGroup {
			return
		}
		if e.Info.Type != "text" {
			return
		}
		text := e.Message.GetConversation()
		if text == "" {
			text = e.Message.GetExtendedTextMessage().GetText()
			if text == "" {
				log.Printf("message conversation is nil %+v", e)
				return
			}
		}
		switch {
		case strings.Contains(text, "Аламан туралы білгім келеді"):
			leadCh <- e.Info.Sender.User
		case strings.Contains(text, "https://www.instagram.com/p/DGYSW_RgJYz/"):
			leadCh <- e.Info.Sender.User
		case strings.Contains(text, "Инстадан көрдім! Аламан өнімдері бойынша"):
			leadCh <- e.Info.Sender.User
		}
	}
}

func ChatEventsHandler(userId int) whHandler {
	return func(client *whatsmeow.Client) func(any) {
		return func(evt any) {
			msg, ok := evt.(*events.Message)
			if !ok {
				return
			}
			if msg.Info.IsGroup {
				return
			}
			msgType := ""
			mediaPath := ""
			text := ""
			if msg.Info.Type == "text" {
				msgType = "text"
				if msg.Message.Conversation != nil {
					text = *msg.Message.Conversation
				}
			} else if msg.Info.Type == "media" {
				mediaType, err := getMsgMediaType(msg.Info.MediaType)
				if err != nil {
					log.Println(err)
					return
				}
				msgType = mediaType
				storedPath, err := storeMedia(msg, client)
				if err != nil {
					log.Println(err)
					return
				}
				mediaPath = storedPath
			} else {
				log.Println(ErrUnsupportedMediaType)
				return
			}
			audioLength := 0
			if msg.Message.AudioMessage != nil {
				audioLength = int(*msg.Message.AudioMessage.Seconds)
			}
			newMsg := Message{
				UserId:      userId,
				Text:        text,
				IsFromUser:  msg.Info.IsFromMe,
				LeadPhone:   msg.Info.Chat.User,
				Type:        msgType,
				Path:        mediaPath,
				AudioLength: audioLength,
			}
			msgCh <- newMsg
			return
		}
	}
}

func getMsgMediaType(mediaType string) (string, error) {
	if mediaType == "ptt" {
		return "audio", nil
	} else if mediaType == "image" {
		return "image", nil
	} else if mediaType == "video" {
		return "video", nil
	}
	return "", ErrUnsupportedMediaType
}

const whMediaPath = "./assets/wh-media"

func storeMedia(evt *events.Message, client *whatsmeow.Client) (string, error) {
	switch {
	case evt.Message.ImageMessage != nil:
		img := evt.Message.GetImageMessage()
		if img != nil {
			data, err := client.Download(img)
			if err != nil {
				return "", fmt.Errorf("Failed to download image: %v", err)
			}
			exts, _ := mime.ExtensionsByType(img.GetMimetype())
			path := fmt.Sprintf("%s/image/%s-%s%s", whMediaPath, evt.Info.Sender.User, evt.Info.ID, exts[0])
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				return "", fmt.Errorf("Failed to save image: %v", err)
			}
			return path, nil
		}
	case evt.Message.AudioMessage != nil:
		audio := evt.Message.GetAudioMessage()
		if audio != nil {
			data, err := client.Download(audio)
			if err != nil {
				return "", fmt.Errorf("Failed to download audio: %v", err)
			}
			exts, _ := mime.ExtensionsByType(audio.GetMimetype())
			path := fmt.Sprintf("%s/audio/%s-%s%s", whMediaPath, evt.Info.Sender.User, evt.Info.ID, exts[0])
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				return "", fmt.Errorf("Failed to save audio: %v", err)
			}
			return path, nil
		}
	case evt.Message.VideoMessage != nil:
		video := evt.Message.GetVideoMessage()
		if video != nil {
			data, err := client.Download(video)
			if err != nil {
				return "", fmt.Errorf("Failed to download video: %v", err)
			}
			exts, _ := mime.ExtensionsByType(video.GetMimetype())
			path := fmt.Sprintf("%s/video/%s-%s%s", whMediaPath, evt.Info.Sender.User, evt.Info.ID, exts[0])
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				return "", fmt.Errorf("Failed to save video: %v", err)
			}
			return path, nil
		}
	}
	return "", ErrUnsupportedMediaType
}
