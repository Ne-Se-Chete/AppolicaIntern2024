package main

import (
	"bytes"
	"container/heap"
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"context"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/sashabaranov/go-openai"
)

// Structs
type User struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type Item struct {
	ID   string `json:"_id"`
	Name string `json:"item name"`
}

type ItemInfo struct {
	ItemName      string `json:"item name"`
	SecondsToCook int    `json:"seconds to cook"`
	CapacityOnGrill int  `json:"capacity on grill"`
}

type Order struct {
	Item     string
	Quantity int
	CookTime int
}

// PriorityQueue implementation
type PriorityQueue []*Order

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Quantity < pq[j].Quantity }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Order)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	order := old[n-1]
	*pq = old[:n-1]
	return order
}

var orderQueue PriorityQueue
var orderDeadline time.Time
var ordersEnabled bool

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
	if err := godotenv.Load(); err != nil {
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
			// Implement interactive event handling if needed
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
	switch cmd.Command {
	case "/hi":
		postMessage(client, cmd.ChannelID, "Hi, I'm Slack Bot. I got your command.")
	case "/order":
		handleOrder(client, cmd)
	case "/start":
		handleStart(client, cmd)
	case "/help":
		message := "This is a Slack bot for managing orders. Here's how it works:\n" +
			"1. Type `/start {time}` to start a new session for orders. The `{time}` argument sets a deadline after which no new orders will be accepted.\n" +
			"2. Type `/order {item_from_the_menu} {quantity}` to place a new order. The `{item_from_the_menu}` argument specifies what you want to eat, and the `quantity` specifies how much you want.\n" +
			"NOTE: You can see the full menu with the command `/menu` and if you want to add a new product, you need to type" + 
			" `/menu add {item} {capacity_on_grill} {price} {seconds_to_cook}` where {item} is the product you want to add, " +
			"{capacity_on_grill} is how many of this items can be placed on the grill at the same type, {price} is how much it costs "+
			"and {seconds_to_cook} is how many seconds it must be cooked (approximately)."
		postMessage(client, cmd.ChannelID, message)
	case "/menu":
		handleMenu(client, cmd)
	case "/receipt":
		handleReceipt(client, cmd)
	default:
		log.Printf("Unknown command: %s", cmd.Command)
	}
}



func handleReceipt(client *slack.Client, cmd slack.SlashCommand) {
	channelID := os.Getenv("CHANNEL_ID")

	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     10,
	}

	history, err := client.GetConversationHistory(&historyParams)
	if err != nil {
		postMessage(client, cmd.ChannelID, "Error getting conversation history.")
		return
	}

	// Find the latest image file
	var latestFile *slack.File
	for _, msg := range history.Messages {
		if len(msg.Files) > 0 {
			file := msg.Files[0] // Get the first file from the message
			if file.Mimetype == "image/png" || file.Mimetype == "image/jpeg" { // Check for image type
				latestFile = &file // Assign the file pointer to latestFile
				break
			}
		}
	}

	if latestFile == nil {
		postMessage(client, cmd.ChannelID, "No image found in the last 10 messages.")
		return
	}

	// Download the image with proper authentication
	req, err := http.NewRequest("GET", latestFile.URLPrivateDownload, nil)
	if err != nil {
		postMessage(client, cmd.ChannelID, "Error creating request.")
		return
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SLACK_BOT_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		postMessage(client, cmd.ChannelID, "Error downloading image.")
		return
	}
	defer resp.Body.Close()

	// Read the image data
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		postMessage(client, cmd.ChannelID, "Error reading image data.")
		return
	}

	// Save the image locally
	fileName := "./images/" + latestFile.Name
	err = ioutil.WriteFile(fileName, imageData, 0644)
	if err != nil {
		postMessage(client, cmd.ChannelID, "Error saving image locally.")
		return
	}

	fmt.Printf("Image saved locally as: %s\n", fileName)
	postMessage(client, cmd.ChannelID, "Image saved locally.")

	description, err := getDescriptionFromAPI(fileName)
	if err != nil {
		postMessage(client, cmd.ChannelID, "Error getting image description from API.")
		return
	}

	postMessage(client, cmd.ChannelID, "Image Description: "+description)
}

func getDescriptionFromAPI(imagePath string) (string, error) {
	apiURL := "http://127.0.0.1:5000/describe-image"

	// Create the JSON payload
	payload := map[string]string{"image_path": imagePath}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON payload: %w", err)
	}

	// Send the POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error sending POST request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the JSON response
	var responseData map[string]string
	err = json.Unmarshal(respBody, &responseData)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling JSON response: %w", err)
	}

	// Check for an error in the response
	if errorMsg, exists := responseData["error"]; exists {
		return "", fmt.Errorf("API error: %s", errorMsg)
	}

	// Return the description from the response
	return responseData["description"], nil
}


func detectImage(imagePath string) (string, error) {
	var model = "gpt-4o"
	var apiKey = os.Getenv("OPENAI_API_KEY")
	var temperature float32 = 0.0

	client := openai.NewClient(apiKey)

	base64Image, err := encodeImage64(imagePath)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: "You are a helpful assistant. Describe the image in one sentence."},
		{Role: "user", Content: fmt.Sprintf(`{"type": "text", "text": "Describe the image."},
                                             {"type": "image_url", "image_url": {"url": "data:image/png;base64,%s"}}`, base64Image)},
	}

	ctx := context.Background()

	response, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
	})

	if err != nil {
		log.Fatalf("ChatCompletion error: %v", err)
	}

	return response.Choices[0].Message.Content, nil
}

func encodeImage64(imagePath string) (string, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	imageData, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imageData), nil
}





func handleStart(client *slack.Client, cmd slack.SlashCommand) {
	args := strings.Fields(cmd.Text)
	if len(args) < 1 {
		postMessage(client, cmd.ChannelID, "Please specify the deadline time (in format HH:MM).")
		return
	}

	timeArg := args[0]
	deadline, err := time.Parse("15:04", timeArg)
	if err != nil {
		postMessage(client, cmd.ChannelID, "Invalid time format. Please use HH:MM format.")
		return
	}

	now := time.Now()
	orderDeadline = time.Date(now.Year(), now.Month(), now.Day(), deadline.Hour(), deadline.Minute(), 0, 0, now.Location())
	ordersEnabled = true
	orderQueue = PriorityQueue{}

	response := fmt.Sprintf("Order session started. You can place orders until %s.", orderDeadline.Format("15:04"))
	postMessage(client, cmd.ChannelID, response)

	go func() {
		time.Sleep(time.Until(orderDeadline))
		summarizeOrders(client, cmd.ChannelID)
	}()
}

func handleOrder(client *slack.Client, cmd slack.SlashCommand) {
	if !ordersEnabled {
		postMessage(client, cmd.ChannelID, "Orders are not enabled. Start a new session with /start {time}.")
		return
	}

	args := strings.Fields(cmd.Text)
	if len(args) < 2 {
		postMessage(client, cmd.ChannelID, "Please specify the item and quantity.")
		return
	}

	item := args[0]
	quantity, err := strconv.Atoi(args[1])
	if err != nil {
		postMessage(client, cmd.ChannelID, "Invalid quantity. Please enter a number.")
		return
	}

	userID := getUserID("Pencho")
	if userID == "" {
		postMessage(client, cmd.ChannelID, "User ID not found.")
		return
	}

	itemID := getItemID(item)
	if itemID == "" {
		postMessage(client, cmd.ChannelID, "Item not found.")
		return
	}

	loc := time.FixedZone("EEST", 2*60*60)
	startTime := time.Now().In(loc).Format("Jan 2, 2006 03:04 pm")

	order := map[string]interface{}{
		"ordered item": itemID,
		"owner":        userID,
		"quantity":     quantity,
		"start time":   startTime,
	}

	sendOrder(order)

	response := fmt.Sprintf("Order placed: %s %d", item, quantity)
	postMessage(client, cmd.ChannelID, response)

	itemData := fetchItemData()
	itemInfo, ok := itemData[item]
	if !ok {
		postMessage(client, cmd.ChannelID, "Failed to fetch item data.")
		return
	}
	cookTime := calculateCookingTime(quantity, itemInfo.CapacityOnGrill, itemInfo.SecondsToCook)
	newOrder := &Order{
		Item:     item,
		Quantity: quantity,
		CookTime: cookTime,
	}

	heap.Push(&orderQueue, newOrder)
}

func fetchItemData() map[string]ItemInfo {
	url := os.Getenv("SERVER_ITEM")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from server: %v", resp.Status)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to parse response body: %v", err)
		return nil
	}

	itemData := make(map[string]ItemInfo)
	if response, ok := result["response"].(map[string]interface{}); ok {
		if results, ok := response["results"].([]interface{}); ok {
			for _, res := range results {
				if record, ok := res.(map[string]interface{}); ok {
					itemName, _ := record["item name"].(string)
					secondsToCook, _ := record["seconds to cook"].(float64)
					capacityOnGrill, _ := record["capacity on grill"].(float64)
					itemData[itemName] = ItemInfo{
						ItemName:        itemName,
						SecondsToCook:   int(secondsToCook),
						CapacityOnGrill: int(capacityOnGrill),
					}
				}
			}
		}
	}

	return itemData
}

func handleMenu(client *slack.Client, cmd slack.SlashCommand) {
	url := os.Getenv("SERVER_ITEM")
		if url == "" {
			log.Printf("SERVER_ITEM environment variable is not set")
			postMessage(client, cmd.ChannelID, "Failed to add the item.")
			return
		}
	args := strings.Fields(cmd.Text)

	if len(args) > 0 && args[0] == "add" && len(args) > 4 {
		// Handle adding an item to the menu

		item := args[1]
		capacityOnGrill := args[2]
		price := args[3]
		secondsToCook := args[4]

		// Create the payload for the POST request
		newItem := map[string]interface{}{
			"item name":        item,
			"capacity on grill": capacityOnGrill,
			"price":            price,
			"seconds to cook":  secondsToCook,
		}

		newItemJSON, err := json.Marshal(newItem)
		if err != nil {
			log.Printf("Failed to marshal new item: %v", err)
			postMessage(client, cmd.ChannelID, "Failed to add the item.")
			return
		}

		req, err := http.NewRequest("POST", url, strings.NewReader(string(newItemJSON)))
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			postMessage(client, cmd.ChannelID, "Failed to add the item.")
			return
		}

		req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))
		req.Header.Add("Content-Type", "application/json")

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf("Failed to send request: %v", err)
			postMessage(client, cmd.ChannelID, "Failed to add the item.")
			return
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusCreated {
			log.Printf("Error response from server: %v", resp.Status)
			postMessage(client, cmd.ChannelID, "Failed to add the item.")
			return
		}

		postMessage(client, cmd.ChannelID, "Successfully added item: "+item)
		return
	}

	// Handle fetching and displaying the menu
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		postMessage(client, cmd.ChannelID, "Failed to fetch the menu.")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from server: %v", resp.Status)
		postMessage(client, cmd.ChannelID, "Failed to fetch the menu.")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		postMessage(client, cmd.ChannelID, "Failed to fetch the menu.")
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to parse response body: %v", err)
		postMessage(client, cmd.ChannelID, "Failed to fetch the menu.")
		return
	}

	var menuItems []string
	if response, ok := result["response"].(map[string]interface{}); ok {
		if results, ok := response["results"].([]interface{}); ok {
			for _, res := range results {
				if record, ok := res.(map[string]interface{}); ok {
					if itemName, ok := record["item name"].(string); ok {
						menuItems = append(menuItems, itemName)
					}
				}
			}
		}
	}

	if len(menuItems) == 0 {
		postMessage(client, cmd.ChannelID, "No items found in the menu.")
		return
	}

	menuMessage := "Here is the menu:\n" + strings.Join(menuItems, "\n")
	postMessage(client, cmd.ChannelID, menuMessage)
}


func calculateCookingTime(quantity, capacity, baseTime int) int {
	batches := quantity / capacity
	if quantity%capacity != 0 {
		batches++
	}
	return batches * baseTime
}

func summarizeOrders(client *slack.Client, channelID string) {
	ordersEnabled = false
	if len(orderQueue) == 0 {
		postMessage(client, channelID, "No orders were placed.")
		return
	}

	itemData := fetchItemData()
	if itemData == nil {
		postMessage(client, channelID, "Failed to fetch item data.")
		return
	}

	orderMap := make(map[string]int)
	for orderQueue.Len() > 0 {
		order := heap.Pop(&orderQueue).(*Order)
		orderMap[order.Item] += order.Quantity
	}

	var summaryBuilder strings.Builder
	summaryBuilder.WriteString("You have collectively ordered:\n")

	counter := 1
	totalCookingTime := 0

	for item, quantity := range orderMap {
		itemInfo, ok := itemData[item]
		if !ok {
			itemInfo = ItemInfo{SecondsToCook: 0, CapacityOnGrill: 1}
		}
		itemCookingTime := calculateCookingTime(quantity, itemInfo.CapacityOnGrill, itemInfo.SecondsToCook)
		totalCookingTime += itemCookingTime
		summaryBuilder.WriteString(fmt.Sprintf("%d. %s x%d -> %d seconds to cook\n", counter, item, quantity, totalCookingTime))
		counter++

		orderSummary := map[string]interface{}{
			"item ordered":    getItemID(item),
			"seconds to cook": totalCookingTime,
			"summed quantity": quantity,
		}
		sendOrderSummary(orderSummary)

	}

	summaryBuilder.WriteString("The order won't be received now - start a new order session with /start {time}.")

	postMessage(client, channelID, summaryBuilder.String())
}

func sendOrderSummary(orderSummary map[string]interface{}) {
	client := &http.Client{}
	url := os.Getenv("SERVER_TODAYS_ORDER")

	orderSummaryJSON, err := json.Marshal(orderSummary)
	if err != nil {
		log.Printf("Failed to marshal order summary: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(orderSummaryJSON)))
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

	log.Println("Order summary sent successfully")
}


func getUserID(userName string) string {
	return getIDFromServer(os.Getenv("SERVER_USERS"), userName, "name")
}

func getItemID(itemName string) string {
	return getIDFromServer(os.Getenv("SERVER_ITEM"), itemName, "item name")
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

func postMessage(client *slack.Client, channelID, message string) {
	if _, _, err := client.PostMessage(channelID, slack.MsgOptionText(message, false)); err != nil {
		log.Printf("Failed to post message: %v", err)
	}
}
