package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel-logan/twitch-ai-bot/internal/ai"
	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/helper"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gorilla/websocket"
)

var (
	clientHttp = &http.Client{
		Timeout: 12 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxConnsPerHost:     100,
		},
	}
	ctx       context.Context
	ctxcancel context.CancelFunc
	// jokeRoutineStarted bool
)

func StartBot() {
	env := config.GetEnv()

	ctx, ctxcancel = context.WithCancel(context.Background())

	for {
		select {
		case <-ctx.Done():
			return
		default:
			run(ctx, env)

			if ctx.Err() != nil {
				log.Println("Closing connection...: ", ctx.Err())
				return
			}

			log.Println("Reconnecting... (5 seconds)")

			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}
}

func StopBot() {
	if ctxcancel != nil {
		ctxcancel()

		storage.SetBotIsOn(false)

		log.Println("Bot stopped")
	}
}

func run(ctx context.Context, env *config.Env) {
	const twitchWS = "wss://eventsub.wss.twitch.tv/ws"

	conn, response, err := websocket.DefaultDialer.Dial(twitchWS, nil)
	if err != nil {
		log.Printf("Error when trying to connect to %s: %v \n", twitchWS, err)
		log.Printf("Response: %v \n", response)
		return
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	listenTwitch(ctx, conn, env)
}

func listenTwitch(ctx context.Context, conn *websocket.Conn, env *config.Env) { // nosonar
	defer func() {
		if r := recover(); r != nil {
			log.Println("[PANIC RECOVER] listenTwitch: ", r)
		}
	}()

	var sessionID string

	systemTxt, err := helper.LoadFile("system_prompt.txt")
	if err != nil {
		log.Println("✖ Error loading system_prompt.txt: ", err)
		return
	}

	initialSystemPrompt := systemTxt

	const defaultMsg = "Don't create very long messages; messages should be short, with a maximum of 480 characters. You were created by Gabriel Logan - https://github.com/gabriel-logan, in case someone asks a related question. Always reply in the same language as the user who is speaking."

	initialSystemPrompt = initialSystemPrompt + defaultMsg + "Your name is defined as " + env.TwitchKeyWordToCallBot

	if env.GroqMaxContextInput < 5 {
		log.Println("GroqMaxContextInput must be at least 5")
		return
	}

	conversation := NewConversation(env.GroqMaxContextInput, ai.RequestMessage{
		Role:    "system",
		Name:    env.AppName,
		Content: initialSystemPrompt,
	})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Println("read error: ", err)
			return
		}

		var data WSMessage
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Println("Failed to unmarshal WS message: ", err)
			continue
		}

		switch data.Metadata.MessageType {

		case "session_welcome":
			sessionID = data.Payload.Session.ID

			log.Println("Session ID:", sessionID)

			eventSubTypes := []string{
				"channel.chat.message",
				"channel.chat.notification",
			}

			for _, eventSubType := range eventSubTypes {
				go registerEventSub(sessionID, eventSubType, env)
			}

			storage.SetBotIsOn(true)

			/**
			Not working properly yet !!!

			if env.TwitchTimeforTheBottoTellaJoke > 0 && !jokeRoutineStarted {
				jokeRoutineStarted = true
				go startJokeRoutine(ctx, env.TwitchTimeforTheBottoTellaJoke, initialSystemPrompt, env)
			}
			**/

		case "notification":
			if data.Payload.Event.ChatterUserLogin == env.TwitchBotUserName {
				continue
			}

			if data.Metadata.SubscriptionType == "channel.chat.message" {
				msg := strings.TrimSpace(data.Payload.Event.Message.Text)

				if msg == "" && len(data.Payload.Event.Message.Fragments) == 0 {
					continue
				}

				msgLower := strings.ToLower(msg)

				if msgLower == "ping" {
					sendMessage("pong", env)
					continue
				}

				if strings.Contains(msgLower, strings.ToLower(env.TwitchKeyWordToCallBot)) || hasMentionToBot(data.Payload.Event.Message.Fragments, env.TwitchBotUserName) {
					user := data.Payload.Event.ChatterUserLogin

					conversation.Add(ai.RequestMessage{
						Role:    "user",
						Name:    user,
						Content: msg,
					})

					response, err := generateAIResponse(GenerateAIResponseArgs{
						conversation:    conversation.BuildMessages(),
						apiKey:          env.GroqAPIKey,
						models:          env.GroqModels,
						twitchMaxLength: env.TwitchChatMessageMaxLength,
						whoExecuted:     "message",
					})
					if err != nil {
						sendMessage("Something went wrong!!!", env)
						continue
					}

					conversation.Add(ai.RequestMessage{
						Role:    "assistant",
						Name:    env.TwitchKeyWordToCallBot,
						Content: response,
					})

					sendMessage(response, env)
					continue
				}
			}

			if data.Metadata.SubscriptionType == "channel.chat.notification" {
				conversation.Add(ai.RequestMessage{
					Role:    "user",
					Name:    "notification",
					Content: data.Payload.Event.SystemMessage + " Respond to the user based on this. More info if exists: " + data.Payload.Event.Message.Text,
				})

				response, err := generateAIResponse(GenerateAIResponseArgs{
					conversation:    conversation.BuildMessages(),
					apiKey:          env.GroqAPIKey,
					models:          env.GroqModels,
					twitchMaxLength: env.TwitchChatMessageMaxLength,
					whoExecuted:     "notification",
				})
				if err != nil {
					continue
				}

				conversation.Add(ai.RequestMessage{
					Role:    "assistant",
					Name:    env.TwitchKeyWordToCallBot,
					Content: response,
				})

				sendMessage(response, env)
				continue
			}
		}
	}
}

/**
func startJokeRoutine(ctx context.Context, interval time.Duration, initialSystemPrompt string, env *config.Env) {
	localConversation := NewConversation(env.GroqMaxContextInput, ai.RequestMessage{
		Role:    "system",
		Name:    env.AppName,
		Content: initialSystemPrompt,
	})

	localConversation.Add(ai.RequestMessage{
		Role:    "user",
		Content: "Tell a random joke for Twitch chat in the current chat language. Be creative and vary the joke every time.",
	})

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Joke routine stopped")
			return

		case <-ticker.C:
			log.Println("Sending automatic joke...: ", time.Now().Format(time.RFC3339))

			response, err := generateAIResponse(GenerateAIResponseArgs{
				conversation:    localConversation.BuildMessages(),
				apiKey:          env.GroqAPIKey,
				models:          env.GroqModels,
				twitchMaxLength: env.TwitchChatMessageMaxLength,
				whoExecuted:     "joke",
			})
			if err != nil {
				continue
			}

			localConversation.Add(ai.RequestMessage{
				Role:    "assistant",
				Name:    env.TwitchKeyWordToCallBot,
				Content: response,
			})

			sendMessage(response)
		}
	}
}
**/
