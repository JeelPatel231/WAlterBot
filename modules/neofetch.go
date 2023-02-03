package modules

import (
	"context"
	"gowa/utils"
	"os/exec"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

const neofetch_help string = "Flex the system specifications."

func neofetch_callback(cli *whatsmeow.Client, msg *events.Message) error {
	stdout, err := exec.Command("neofetch", "--stdout").Output()
	if err != nil {
		return err
	}
	_, err = cli.SendMessage(
		context.Background(),
		msg.Info.Chat,
		utils.NewMessage(string(stdout), msg),
	)

	return err
}

var Neofetch Command = Command{
	".neofetch",
	neofetch_callback,
	neofetch_help,
}
