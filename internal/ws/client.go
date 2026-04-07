package ws

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
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

var once sync.Once

func StartBot() {
	once.Do(func() {
		go func() {
			for {
				Run()
				log.Println("Reconnecting... (5 seconds)")
				time.Sleep(5 * time.Second)
			}
		}()
	})
}

func Run() {
	const twitchWS = "wss://eventsub.wss.twitch.tv/ws"

	conn, response, err := websocket.DefaultDialer.Dial(twitchWS, nil)
	if err != nil {
		log.Printf("Error when trying to connect to %s: %v \n", twitchWS, err)
		log.Printf("Response: %v \n", response)
		return
	}
	defer conn.Close()

	ListenTwitch(conn, config.GetEnv())
}

func ListenTwitch(conn *websocket.Conn, env *config.Env) { // nosonar
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
				if strings.TrimSpace(data.Payload.Event.Message.Text) == "ping" {
					go sendMessage(env, "pong")
				}

				if strings.Contains(strings.ToLower(data.Payload.Event.Message.Text), "jesus") {
					go sendMessage(env, "Jesus é o Senhor")
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
}
