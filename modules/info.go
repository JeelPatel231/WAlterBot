package modules

import (
	"context"
	"gowa/utils"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

const info_help string = "Get info of the replied entity"

func info_callback(cli *whatsmeow.Client, msg *events.Message) {
	if msg.Message.ExtendedTextMessage == nil || msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage == nil {
		cli.SendMessage(
			context.Background(),
			msg.Info.Chat,
			utils.NewMessage("Reply to a user's message dumbass", msg),
		)
		return
	}

	cli.SendMessage(
		context.Background(),
		msg.Info.Chat,
		utils.NewMessage("Not Implemented yet! :P", msg),
	)
}

var Info Command = Command{
	".info",
	info_callback,
	info_help,
}
