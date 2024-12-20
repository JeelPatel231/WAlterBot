package modules

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
	Command  string
	Callback func(client *whatsmeow.Client, message *events.Message) error
	Help     string
}

var CommandArray = []Command{
	Hello,
	Sticker,
	Info,
	Neofetch,
	Ytdl,
}
