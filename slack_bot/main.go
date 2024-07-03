package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type User struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type Item struct {
	ID   string `json:"_id"`
	Name string `json:"item name"`
}

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
	args := strings.Fields(cmd.Text)
	if len(args) < 2 {
		_, _, err := client.PostMessage(cmd.ChannelID, slack.MsgOptionText("Please specify the item and quantity.", false))
		if err != nil {
			log.Printf("Failed to post message: %v", err)
		}
		return
	}

	item := args[0]
	quantity := args[1]

	userID := getUserID("Pencho")
	if userID == "" {
		_, _, err := client.PostMessage(cmd.ChannelID, slack.MsgOptionText("User ID not found.", false))
		if err != nil {
			log.Printf("Failed to post message: %v", err)
		}
		return
	}

	itemID := getItemID(item)
	if itemID == "" {
		_, _, err := client.PostMessage(cmd.ChannelID, slack.MsgOptionText("Item not found.", false))
		if err != nil {
			log.Printf("Failed to post message: %v", err)
		}
		return
	}

	loc := time.FixedZone("EEST", 2*30)
	startTime := time.Now().In(loc).Format("Jan 2, 2006 03:04 pm")

	order := map[string]interface{}{
		"ordered item":   itemID,
		"owner":   userID,
		"quantity":  quantity,
		"start time": startTime,
	}

	sendOrder(order)
	response := fmt.Sprintf("Order placed: %s %s", item, quantity)
	_, _, err := client.PostMessage(cmd.ChannelID, slack.MsgOptionText(response, false))
	if err != nil {
		log.Printf("Failed to post message: %v", err)
	}
}

func getUserID(userName string) string {
	url := os.Getenv("SERVER_USERS")
	return getIDFromServer(url, userName, "name")
}

func getItemID(itemName string) string {
	url := os.Getenv("SERVER_ITEM")
	return getIDFromServer(url, itemName, "item name")
}

func getIDFromServer(url, searchValue, searchKey string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return ""
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from server: %v", resp.Status)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to parse response body: %v", err)
		return ""
	}

	if response, ok := result["response"].(map[string]interface{}); ok {
		if results, ok := response["results"].([]interface{}); ok {
			for _, res := range results {
				if record, ok := res.(map[string]interface{}); ok {
					if name, ok := record[searchKey].(string); ok && name == searchValue {
						if id, ok := record["_id"].(string); ok {
							return id
						}
					}
				}
			}
		}
	}

	log.Println("ID not found in response")
	return ""
}

func sendOrder(order map[string]interface{}) {
	client := &http.Client{}
	url := os.Getenv("SERVER_ORDER")

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshal order: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(orderJSON)))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from server: %v", resp.Status)
		return
	}

	log.Println("Order sent successfully")
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
