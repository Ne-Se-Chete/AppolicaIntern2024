package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	loadEnv()

	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")

	validateTokens(appToken, botToken)

	client := createSlackClient(botToken, appToken)
	socketClient := createSocketClient(client)

	go handleEvents(socketClient, client)

	socketClient.Run()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func validateTokens(appToken, botToken string) {
	if appToken == "" || botToken == "" {
		log.Fatalf("SLACK_APP_TOKEN and SLACK_BOT_TOKEN must be set")
	}
}

func createSlackClient(botToken, appToken string) *slack.Client {
	return slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(appToken),
	)
}

func createSocketClient(client *slack.Client) *socketmode.Client {
	return socketmode.New(
		client,
		socketmode.OptionDebug(true),
	)
}

func handleEvents(socketClient *socketmode.Client, client *slack.Client) {
	for evt := range socketClient.Events {
		switch evt.Type {
		case socketmode.EventTypeInteractive:
			handleInteractiveMessage(client, evt.Data.(slack.InteractionCallback))
		case socketmode.EventTypeSlashCommand:
			cmd, ok := evt.Data.(slack.SlashCommand)
			if !ok {
				log.Printf("Ignored %+v\n", evt)
				continue
			}
			socketClient.Ack(*evt.Request)
			handleSlashCommand(client, cmd)
		default:
			log.Printf("Ignored %+v\n", evt)
		}
	}
}

func handleSlashCommand(client *slack.Client, cmd slack.SlashCommand) {
	switch {
	case cmd.Command == "/hi":
		_, _, err := client.PostMessage(cmd.ChannelID, slack.MsgOptionText("Hi, I'm Test Bot. I got your command.", false))
		if err != nil {
			log.Printf("Failed to post message: %v", err)
		}
	case cmd.Command == "/order":
		handleOrder(client, cmd)
	default:
		log.Printf("Unknown command: %s", cmd.Command)
	}
}

func handleOrder(client *slack.Client, cmd slack.SlashCommand) {
	order := strings.TrimSpace(strings.TrimPrefix(cmd.Text, "order"))

	if strings.Contains(strings.ToLower(order), "kufte") {
		sendKufteOrderMessage(client, cmd.ChannelID)
	} else {
		response := "Received order: " + order
		_, _, err := client.PostMessage(cmd.ChannelID, slack.MsgOptionText(response, false))
		if err != nil {
			log.Printf("Failed to post message: %v", err)
		}
	}
}

func sendKufteOrderMessage(client *slack.Client, channelID string) {
	attachment := slack.Attachment{
		Text:       "Please choose the size of kufte:",
		CallbackID: "kufte_order",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			{
				Name:  "kufte_size",
				Type:  "select",
				Text:  "Select size",
				Options: []slack.AttachmentActionOption{
					{Text: "Small", Value: "small"},
					{Text: "Big", Value: "big"},
				},
			},
		},
	}

	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionText("Please select the size of kufte:", false),
	)
	if err != nil {
		log.Printf("Failed to post message: %v", err)
	}
}

func handleInteractiveMessage(client *slack.Client, interaction slack.InteractionCallback) {
	switch interaction.CallbackID {
	case "kufte_order":
		if len(interaction.ActionCallback.BlockActions) > 0 {
			selectedOption := interaction.ActionCallback.BlockActions[0].SelectedOption.Value
			order := "kufte " + selectedOption
			response := "Order - " + order
			_, _, err := client.PostMessage(interaction.Channel.ID, slack.MsgOptionText(response, false))
			if err != nil {
				log.Printf("Failed to post message: %v", err)
				return
			}
			log.Printf("Message sent successfully.")
		} else {
			log.Printf("No block actions found in interaction callback.")
		}
	default:
		log.Printf("Unhandled interactive message callback: %s", interaction.CallbackID)
	}
}
