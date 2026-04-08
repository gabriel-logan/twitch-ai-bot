package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gabriel-logan/twitch-ai-bot/internal/ws"
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Data []struct {
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
	} `json:"data"`
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
	q.Add("login", c.Query("login"))

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

func StartTwitchBot(c *gin.Context) {
	isBotOn := storage.GetBotIsOn()

	if isBotOn {
		c.JSON(http.StatusOK, "Bot already started")
		return
	}

	ws.StartBot()

	c.JSON(http.StatusOK, "Bot started")
}
