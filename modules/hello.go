package modules

import (
	"context"
	"gowa/utils"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

const hello_help string = "Check if Bot is alive or dead."

func hello_callback(cli *whatsmeow.Client, msg *events.Message) error {
	_, err := cli.SendMessage(
		context.Background(),
		msg.Info.Chat,
		utils.NewMessage("Hello World", msg),
	)
	return err
}

var Hello Command = Command{
	".hello",
	hello_callback,
	hello_help,
}
