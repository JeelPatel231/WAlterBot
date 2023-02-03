package modules

import (
	"context"
	"errors"
	"gowa/utils"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

const info_help string = "Get info of the replied entity"

func info_callback(cli *whatsmeow.Client, msg *events.Message) error {
	if msg.Message.ExtendedTextMessage == nil || msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage == nil {
		return errors.New("Reply to a user's message to get info!")
	}

	cli.SendMessage(
		context.Background(),
		msg.Info.Chat,
		utils.NewMessage("Not Implemented yet! :P", msg),
	)
	return nil
}

var Info Command = Command{
	".info",
	info_callback,
	info_help,
}
