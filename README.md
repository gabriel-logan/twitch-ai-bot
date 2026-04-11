# Twitch AI Bot

A Twitch chatbot powered by AI that listens to your channel's chat and responds intelligently when invoked by a keyword. Supports multiple AI providers (cloud) with automatic fallback. Built with Go, Gin, and Twitch EventSub WebSocket API.

## Features

- **Multi-Provider AI Support**: Integrates with multiple AI providers (cloud APIs) with automatic fallback if one goes down
- **AI-Powered Responses**: Generates intelligent, context-aware responses in chat
- **Twitch EventSub Integration**: Uses Twitch's WebSocket EventSub for real-time chat monitoring
- **OAuth2 Authentication**: Secure Twitch login flow with proper scopes (`user:read:chat`, `user:write:chat`, `user:bot`, `channel:bot`)
- **Conversation Context**: Maintains a history of the entire conversation context, with a limit defined by the user in the environment variable
- **Web Dashboard**: Simple UI for authentication management and Twitch user info lookup
- **Auto-Reconnection**: Automatically reconnects to Twitch WebSocket on disconnection
- **Custom System Prompt**: Configurable system prompt via `system_prompt.txt`
- **Keyword Activation**: Bot only responds when a configurable keyword is mentioned (default: `jesus`)

## Prerequisites

- Go 1.25 or higher
- A Twitch Developer Application ([create one here](https://dev.twitch.tv/console))
- A Groq API key ([get one here](https://console.groq.com/keys))
- A Twitch channel where the bot will operate

## Installation

1. **Clone the repository**:

```bash
git clone https://github.com/gabriel-logan/twitch-ai-bot.git
cd twitch-ai-bot
```

2. **Install dependencies**:

```bash
go mod download
```

3. **Configure environment variables**:

```bash
cp .env.example .env
```

Edit `.env` with your credentials (see [Configuration](#configuration) below).

4. **Run the application**:

```bash
make run
```

Or build a binary:

```bash
make build
./bin/twitch_ai_bot
```

## Configuration

All configuration is done via the `.env` file:

| Variable | Description | Example |
|---|---|---|
| `GIN_MODE` | Gin server mode (`debug` or `release`) | `debug` |
| `APP_NAME` | Application display name | `TwitchAIBot` |
| `SERVER_PORT` | HTTP server port | `8080` |
| `SERVER_TRUSTED_PROXIES` | List of trusted proxy addresses for the server | `10.0.0.0/8,172.16.0.0/12,172.12.2.1` |
| `TWITCH_CLIENT_ID` | Twitch OAuth Client ID | `abc123...` |
| `TWITCH_CLIENT_SECRET` | Twitch OAuth Client Secret | `xyz789...` |
| `TWITCH_CLIENT_REDIRECT_URI` | OAuth callback URL | `http://localhost:8080/api/auth/callback/twitch` |
| `TWITCH_BROADCASTER_ID` | Twitch ID of the channel to monitor | `12345678` |
| `TWITCH_BOT_USER_ID` | Twitch ID of the bot account | `87654321` |
| `TWITCH_BOT_USER_NAME` | Login name of the bot account | `my_bot_user` |
| `TWITCH_KEY_WORD_TO_CALL_BOT` | Keyword that triggers the bot | `jesus` |
| `TWITCH_CHAT_MESSAGE_MAX_LENGTH` | Twitch chat message max length | `500` |
| `GROQ_API_KEY` | Groq API key | `gsk_...` |
| `GROQ_MODELS` | Groq models to use - automatic fallback | `llama-3.3-70b-versatile,openai/gpt-oss-120b` |
| `GROQ_MAX_CONTEXT_INPUT` | Max messages kept in conversation context - minimum is 5 | `25` |
| `CONTEXT_REQUEST_DURATION` | Maximum request timeout | `15s` |

### Setting Up Twitch OAuth

1. Go to [Twitch Developer Console](https://dev.twitch.tv/console)
2. Click **Register Your Application**
3. Set the **OAuth Redirect URL** to match your `TWITCH_CLIENT_REDIRECT_URI`
4. Copy the **Client ID** and generate a **Client Secret**
5. Add them to your `.env` file

### Getting Twitch User IDs

You can find Twitch User IDs using tools like [Twitch ID Finder](https://tools.streamscharts.com/twitch-id-finder) or via the Twitch API.

## Usage

### 1. Start the Server

```bash
make run
```

The server will start at `http://localhost:8080`.

### 2. Authenticate with Twitch

1. Open `http://localhost:8080` in your browser
2. Click **Login with Twitch**
3. Authorize the application with the requested scopes
4. After successful authentication, the bot will automatically start listening to chat

### 3. Interact with the Bot

In your Twitch channel chat, mention the configured keyword (default: `jesus`) followed by your message:

```
jesus hello, how are you?
jesus what do you think about Go?
```

The bot will respond in the same language as the user.

### 4. Web Dashboard

- **Logged out**: Shows a login button to authenticate with Twitch
- **Logged in**: Shows a logout button and a user info lookup tool to fetch Twitch user details by login name

## How It Works

1. **Authentication**: User logs in via Twitch OAuth2 flow, receiving an access token with chat read/write permissions
2. **WebSocket Connection**: After authentication, the bot connects to `wss://eventsub.wss.twitch.tv/ws`
3. **EventSub Subscription**: The bot subscribes to `channel.chat.message` events for the specified broadcaster
4. **Message Processing**: When a chat message contains the trigger keyword, the bot:
   - Builds a conversation context for that user (up to `GROQ_MAX_CONTEXT_INPUT` messages)
   - Sends the conversation to the configured AI provider
   - Posts the AI response back to the Twitch chat
5. **Auto-Reconnect**: If the WebSocket disconnects, the bot automatically attempts to reconnect every 5 seconds

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Author

Created by [Gabriel Logan](https://github.com/gabriel-logan)
