package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gabriel-logan/twitch-ai-bot/internal/ws"
	"github.com/gin-gonic/gin"
)

type UserInfoData struct {
	BroadcasterType string `json:"broadcaster_type"`
	CreatedAt       string `json:"created_at"`
	Description     string `json:"description"`
	DisplayName     string `json:"display_name"`
	ID              string `json:"id"`
	Login           string `json:"login"`
	OfflineImageURL string `json:"offline_image_url"`
	ProfileImageURL string `json:"profile_image_url"`
	Type            string `json:"type"`
	ViewCount       int    `json:"view_count"`
}

type UserInfo struct {
	Data []UserInfoData `json:"data"`
}

type Environment struct {
	TwitchBroadcasterID string `json:"twitch_broadcaster_id"`
	TwitchBotUserID     string `json:"twitch_bot_user_id"`
	TwitchBotUserName   string `json:"twitch_bot_user_name"`
	TwitchKeyWordToCall string `json:"twitch_key_word_to_call"`
}

func GetTwitchUserInfo(c *gin.Context) {
	env := config.GetEnv()

	ctx, cancel := context.WithTimeout(c.Request.Context(), env.ContextRequestDuration)
	defer cancel()

	const baseURL = "https://api.twitch.tv/helix/users"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		log.Printf("Error when trying to create request: %v \n", err)
		c.JSON(http.StatusInternalServerError, "Error when trying to create request")
		return
	}

	q := req.URL.Query()

	logins := c.QueryArray("login")

	if len(logins) == 0 {
		c.JSON(http.StatusBadRequest, "No logins provided, please provide at least one login")
		return
	}

	for _, login := range logins {
		if login == "" {
			continue
		}

		q.Add("login", login)
	}

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Client-Id", env.TwitchClientID)
	req.Header.Set("Authorization", "Bearer "+storage.GetOauthToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error when trying to create request: %v \n", err)
		c.JSON(http.StatusInternalServerError, "Error when trying to create request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error when trying to read response body: %v \n", err)
		c.JSON(http.StatusInternalServerError, "Error when trying to read response body")
		return
	}

	var marshaledBody UserInfo
	if err := json.Unmarshal(body, &marshaledBody); err != nil {
		log.Printf("Error when trying to unmarshal response body: %v \n", err)
		c.JSON(http.StatusInternalServerError, "Error when trying to unmarshal response body")
		return
	}

	c.JSON(http.StatusOK, marshaledBody)
}

func SetEnvironment(c *gin.Context) {
	env := Environment{
		TwitchBroadcasterID: c.Query("twitch_broadcaster_id"),
		TwitchBotUserID:     c.Query("twitch_bot_user_id"),
		TwitchBotUserName:   c.Query("twitch_bot_user_name"),
		TwitchKeyWordToCall: c.Query("twitch_key_word_to_call"),
	}

	if env.TwitchBroadcasterID != "" {
		os.Setenv("TWITCH_BROADCASTER_ID", env.TwitchBroadcasterID)
	}

	if env.TwitchBotUserID != "" {
		os.Setenv("TWITCH_BOT_USER_ID", env.TwitchBotUserID)
	}

	if env.TwitchBotUserName != "" {
		os.Setenv("TWITCH_BOT_USER_NAME", env.TwitchBotUserName)
	}

	if env.TwitchKeyWordToCall != "" {
		os.Setenv("TWITCH_KEY_WORD_TO_CALL_BOT", env.TwitchKeyWordToCall)
	}

	config.ReloadEnv()

	c.JSON(http.StatusOK, "Environment set")
}

func StartTwitchBot(c *gin.Context) {
	isBotOn := storage.GetBotIsOn()

	if isBotOn {
		c.JSON(http.StatusOK, "Bot already started")
		return
	}

	go ws.StartBot()

	c.JSON(http.StatusOK, "Bot started")
}

func StopTwitchBot(c *gin.Context) {
	isBotOn := storage.GetBotIsOn()

	if !isBotOn {
		c.JSON(http.StatusOK, "Bot already stopped")
		return
	}

	go ws.StopBot()

	c.JSON(http.StatusOK, "Bot stopped")
}
