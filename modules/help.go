package modules

import (
	"context"
	"fmt"
	"gowa/utils"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

const help_help string = "show help for modules, if defined."

func help_callback(cli *whatsmeow.Client, msg *events.Message) {
	text := utils.GetText(msg)
	if text == nil {
		fmt.Println("No text found????")
		return
	}

	stringArr := strings.Fields(*text)

	if len(stringArr) != 2 {
		fmt.Println("malformed input")
		fmt.Println(len(stringArr))
		fmt.Println(stringArr)
		return
	}

	switch stringArr[1] {
	case ".help":
		{

			cli.SendMessage(
				context.Background(),
				msg.Info.Chat,
				utils.NewMessage(help_help, msg),
			)
			return
		}
	case "all":
		{
			help_str := "Available Modules :\n"
			for i := range commandArray {
				help_str += "- " + commandArray[i].command + "\n"
			}
			cli.SendMessage(
				context.Background(),
				msg.Info.Chat,
				utils.NewMessage(help_str, msg),
			)
			return
		}
	}

	for i := range commandArray {
		if commandArray[i].command == stringArr[1] {
			cli.SendMessage(
				context.Background(),
				msg.Info.Chat,
				utils.NewMessage(commandArray[i].help, msg),
			)
		}
	}
}

var Help Command = Command{
	".help",
	help_callback,
	help_help,
}
