package modules

import (
	"context"
	"errors"
	"gowa/utils"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

const help_help string = "show help for modules, if defined."

func help_callback(cli *whatsmeow.Client, msg *events.Message) error {
	text := utils.GetText(msg)
	if text == nil {
		return errors.New("No text found???")
	}

	stringArr := strings.Fields(*text)

	if len(stringArr) != 2 {
		return errors.New("Malformed input")
	}

	switch stringArr[1] {
	case ".help":
		{
			_, err := cli.SendMessage(
				context.Background(),
				msg.Info.Chat,
				utils.NewMessage(help_help, msg),
			)
			return err
		}
	case "all":
		{
			help_str := "Available Modules :\n"
			for i := range commandArray {
				help_str += "- " + commandArray[i].command + "\n"
			}
			_, err := cli.SendMessage(
				context.Background(),
				msg.Info.Chat,
				utils.NewMessage(help_str, msg),
			)
			return err
		}
	}

	for i := range commandArray {
		if commandArray[i].command == stringArr[1] {
			_, err := cli.SendMessage(
				context.Background(),
				msg.Info.Chat,
				utils.NewMessage(commandArray[i].help, msg),
			)
			return err
		}
	}
	return nil
}

var Help Command = Command{
	".help",
	help_callback,
	help_help,
}
