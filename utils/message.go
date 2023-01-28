package utils

import (
	"strings"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func NewMessage(text string, replied *events.Message) *waProto.Message {
	trimmedText := strings.TrimSpace(text)
	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: &trimmedText,
		},
	}

	if replied != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      &replied.Info.ID,
			Participant:   proto.String(replied.Info.MessageSource.Sender.String()),
			QuotedMessage: replied.Message,
		}
	}

	return msg
}

func GetText(message *events.Message) *string {
	if message.Message.Conversation != nil {
		return message.Message.Conversation
	}
	if message.Message.ExtendedTextMessage != nil {
		return message.Message.ExtendedTextMessage.Text
	}
	return nil
}

func Contains[K comparable](arr []K, el K) bool {
	for i := range arr {
		if arr[i] == el {
			return true
		}
	}
	return false
}
