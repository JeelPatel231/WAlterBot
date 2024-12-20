package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"gowa/modules"
	"gowa/utils"
)

type MyClient struct {
	WAClient *whatsmeow.Client
}

func (mycli *MyClient) registerCommand(cmd *modules.Command) {
	// build a command Event handler and add that to the client
	commandEvtHandler := mycli.buildCommandEventHandler(cmd)
	fmt.Println(cmd.Command, "commmand registered")

	mycli.WAClient.AddEventHandler(commandEvtHandler)
}

var allowedUsers []string = []string{
	// ALLOW OTHER USERS TO ACCESS YOUR COMMANDS
	// "91XXXXXXXXXX",
	// ^ Input format
}

func (mycli *MyClient) buildCommandEventHandler(cmd *modules.Command) whatsmeow.EventHandler {
	return func(evt interface{}) {
		// Handle event and access mycli.WAClient
		switch v := evt.(type) {
		case *events.Message:
			isAllowed := utils.Contains(allowedUsers, v.Info.Sender.User)
			if v.Info.IsFromMe || isAllowed {
				// get the command
				msg_text := utils.GetText(v)
				if msg_text == nil {
					// return if no text
					return
				}
				cmd_string := strings.Fields(*msg_text)
				if len(cmd_string) == 0 {
					// return if no command
					return
				}
				if cmd_string[0] != cmd.Command {
					// return if command is not the one we are looking for
					return
				}
				// execute the command callback
				go cmd.Callback(mycli.WAClient, v)
			}
		}
	}
}

func main() {

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	customClient := MyClient{client}

	// register all the command event handlers.
	for i := range modules.CommandArray {
		customClient.registerCommand(&modules.CommandArray[i])
	}

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
