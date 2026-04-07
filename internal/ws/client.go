package ws

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/groq"
	"github.com/gabriel-logan/twitch-ai-bot/internal/helper"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Metadata struct {
		MessageType      string `json:"message_type"`
		SubscriptionType string `json:"subscription_type"`
	} `json:"metadata"`
	Payload struct {
		Session struct {
			ID string `json:"id"`
		} `json:"session"`
		Event struct {
			BroadcasterUserLogin string `json:"broadcaster_user_login"`
			ChatterUserLogin     string `json:"chatter_user_login"`
			Message              struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"event"`
	} `json:"payload"`
}

var (
	conversations = make(map[string][]groq.Message)
	once          sync.Once
)

func StartBot() {
	once.Do(func() {
		go func() {
			for {
				run()
				log.Println("Reconnecting... (5 seconds)")
				time.Sleep(5 * time.Second)
			}
		}()
	})
}

func run() {
	const twitchWS = "wss://eventsub.wss.twitch.tv/ws"

	conn, response, err := websocket.DefaultDialer.Dial(twitchWS, nil)
	if err != nil {
		log.Printf("Error when trying to connect to %s: %v \n", twitchWS, err)
		log.Printf("Response: %v \n", response)
		return
	}
	defer conn.Close()

	listenTwitch(conn, config.GetEnv())
}

func listenTwitch(conn *websocket.Conn, env *config.Env) { // nosonar
	var sessionID string

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}

		var data WSMessage
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Println("json error:", err)
			continue
		}

		switch data.Metadata.MessageType {

		case "session_welcome":
			sessionID = data.Payload.Session.ID

			log.Println("Session ID:", sessionID)

			go registerEventSub(sessionID, env)

		case "notification":
			if data.Metadata.SubscriptionType == "channel.chat.message" {
				if data.Payload.Event.ChatterUserLogin == env.TwitchBotUserName {
					continue
				}

				msg := strings.ToLower(strings.TrimSpace(data.Payload.Event.Message.Text))

				if msg == "ping" {
					sendMessage(env, "pong")
				}

				if strings.Contains(msg, env.TwitchKeyWordToCallBot) {
					user := data.Payload.Event.ChatterUserLogin

					if _, exist := conversations[user]; !exist {
						systemTxt, err := helper.LoadFile("system_prompt.txt")
						if err != nil {
							log.Println("✖ Error loading system_prompt.txt:", err)
							log.Println("⚠ Even if you don't want custom settings, create a system_prompt.txt file at the same level as the groq executable. It can be empty, but it needs to exist. This is required for the program to run.")
							return
						}

						initialSystemPrompt := systemTxt

						initialSystemPrompt = initialSystemPrompt + "Do not create very long messages; messages should be short, with a maximum of 400 characters. You were created by Gabriel Logan; you're a bot designed to help people watching Say Seven's live streams on Twitch. Always respond in the same language as the user speaking."

						conversations[user] = []groq.Message{
							{
								Role:    "system",
								Content: initialSystemPrompt,
							},
						}
					}

					conversations[user] = append(conversations[user], groq.Message{
						Role:    "user",
						Content: msg,
					})

					response, err := groq.CallGroq(conversations[user])
					if err != nil {
						log.Println("groq error:", err)
						continue
					}

					conversations[user] = append(conversations[user], groq.Message{
						Role:    "assistant",
						Content: response,
					})

					maxMessages := env.GroqMaxContextInput

					if len(conversations[user]) > maxMessages {
						conversations[user] = append(conversations[user][:1], conversations[user][len(conversations[user])-maxMessages:]...)
					}

					sendMessage(env, response)
				}
			}
		}
	}
}

func registerEventSub(sessionID string, env *config.Env) {
	token := storage.GetOauthToken()

	body := map[string]interface{}{
		"type":    "channel.chat.message",
		"version": "1",
		"condition": map[string]string{
			"broadcaster_user_id": env.TwitchBroadcasterID,
			"user_id":             env.TwitchBotUserID,
		},
		"transport": map[string]string{
			"method":     "websocket",
			"session_id": sessionID,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Println("json error:", err)
		return
	}

	const baseURL = "https://api.twitch.tv/helix/eventsub/subscriptions"

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("eventsub error:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", env.TwitchClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("eventsub error:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("EventSub status:", resp.Status)
}

func sendMessage(env *config.Env, message string) {
	token := storage.GetOauthToken()

	body := map[string]string{
		"broadcaster_id": env.TwitchBroadcasterID,
		"sender_id":      env.TwitchBotUserID,
		"message":        message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Println("json error:", err)
		return
	}

	const baseURL = "https://api.twitch.tv/helix/chat/messages"

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("send message error:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", env.TwitchClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("send message error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Println("send message failed:", resp.Status)

		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Println("response body:", string(bodyBytes))

		return
	}
}
