package modules

import (
	"context"
	"fmt"
	"gowa/utils"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
	command  string
	callback func(client *whatsmeow.Client, message *events.Message) error
	help     string
}

var commandArray = []Command{
	Hello,
	Sticker,
	Info,
	Neofetch,
	Ytdl,
}

func CallbackExecutor(client *whatsmeow.Client, message *events.Message) {
	// get the command
	msg_text := utils.GetText(message)
	if msg_text == nil {
		return
	}
	cmd_string := strings.Fields(*msg_text)
	if len(cmd_string) == 0 {
		return
	}

	fmt.Println(message.Message.String())

	// show help text if help is called
	if cmd_string[0] == Help.command {
		go Help.callback(client, message)
		return
	}

	// or find the callback and execute it from the list of registered commands
	for i := range commandArray {
		// fmt.Println(commandArray[i], cmd_string[0])
		if commandArray[i].command == cmd_string[0] {
			fmt.Println()
			fmt.Println("Command Captured :", commandArray[i].command)
			fmt.Println("on message:", message.Message.String())
			fmt.Println()
			go func() {
				fmt.Println(commandArray[i])
				err := commandArray[i].callback(client, message)
				if err != nil {
					client.SendMessage(context.Background(), message.Info.Chat,
						utils.NewMessage(err.Error(), message),
					)
				}
			}()
			break
		}
	}

}
