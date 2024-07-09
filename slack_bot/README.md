# Slack Order Management Bot

## What is it used for?

This bot is designed to manage and track orders in a Slack workspace. It allows users to:
- Start a new session for placing orders.
- Place orders specifying items and quantities.
- View the menu and add new items to it.
- Receive notifications about the order deadline.
- Summarize and display the total orders placed.

## How to Set It Up

### Installation of Go

1. **Download and install Go:**
    - Follow the instructions for your operating system on the [official Go installation page](https://golang.org/doc/install).

2. **Verify installation:**
    - Open a terminal and type:
      ```sh
      go version
      ```
    - You should see the installed version of Go.

### Slack API Setup

1. **Create a Slack App:**
    - Go to the [Slack API](https://api.slack.com/apps) and click "Create New App".
    - Choose "From scratch", give your app a name, and select the workspace where you want to install the app.

2. **Add required scopes:**
    - Navigate to "OAuth & Permissions" and add the following bot token scopes:
      - `channels:history`
      - `channels:read`
      - `chat:write`
      - `commands`
      - `files:read`
    - Install the app to your workspace and note down the **Bot User OAuth Token** and **App-Level Token**.

3. **Create Slash Commands:**
    - Go to "Slash Commands" in your Slack app settings.
    - Create commands like `/hi`, `/order`, `/start`, `/help`, `/menu`, and `/receipt`.
    - Set the request URL to the endpoint where your bot will be running.

4. **Set environment variables:**
    - Create a `.env` file in your project directory and add:
      ```env
      SLACK_APP_TOKEN=your-app-level-token
      SLACK_BOT_TOKEN=your-bot-user-oauth-token
      CHANNEL_ID=your-channel-id
      SERVER_ITEM=your-server-item-url
      SERVER_ORDER=your-server-order-url
      SERVER_TODAYS_ORDER=your-server-todays-order-url
      SERVER_USERS=your-server-users-url
      BEARER_TOKEN=your-bearer-token
      OPENAI_API_KEY=your-openai-api-key
      ```

### Running the Bot

1. **Install dependencies:**
    - Run the following command to install required Go packages:
      ```sh
      go get github.com/joho/godotenv github.com/slack-go/slack github.com/slack-go/slack/socketmode github.com/sashabaranov/go-openai
      ```

2. **Build and run the bot:**
    - In your project directory, build the bot:
      ```sh
      go build -o slack-bot
      ```
    - Run the bot:
      ```sh
      ./slack-bot
      ```

## How it Works

### Commands

- **`/hi`**:
    - Responds with a simple greeting message.

- **`/order {item} {quantity}`**:
    - Places a new order for the specified item and quantity.
    - Example: `/order burger 2`
    - Note: This command only adds predefined items that are retrieved from a database. Use this command after the `/start` command.

- **`/start {time}`**:
    - Starts a new session for orders with a deadline.
    - `{time}` should be in `HH:MM` format.
    - Example: `/start 18:30`

- **`/help`**:
    - Displays help information about using the bot.

- **`/menu`**:
    - Displays the current menu.
    - Use `/menu add {item} {capacity_on_grill} {price} {seconds_to_cook}` to add a new item to the menu.
    - Example: `/menu add burger 4 5.99 300`

- **`/receipt`**:
    - Fetches and describes the latest receipt from the Slack channel history with name "receipt".

### Order Workflow

1. **Starting a session**:
    - Use `/start {time}` to initiate an order session with a specific deadline.

2. **Placing orders**:
    - Users place orders using `/order {item} {quantity}`. The bot retrieves the list of available items from the database.

3. **Receiving notifications**:
    - The bot sends reminders about the order deadline, including a 5-minute warning.

4. **Summarizing orders**:
    - Once the deadline is reached, the bot summarizes the orders and posts the total quantities and estimated cooking time.

### Menu Management

- **Fetching the menu**:
    - The bot retrieves the menu from the database.

- **Adding new items**:
    - Users can add new items to the menu using the `/menu add` command, specifying the item details. The new items are stored in the database.

By following the above steps, you can set up and use the Slack Order Management Bot in your workspace effectively.