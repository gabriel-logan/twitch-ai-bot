package handler

import (
	"net/http"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	env := config.GetEnv()

	var IsLoggedIn bool
	if storage.GetOauthToken() != "" {
		IsLoggedIn = true
	}

	isBotOn := storage.GetBotIsOn()

	currentTwitchBroadcasterID := env.TwitchBroadcasterID
	currentTwitchBotUserID := env.TwitchBotUserID
	currentTwitchBotUserName := env.TwitchBotUserName

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      env.AppName,
		"isLoggedIn": IsLoggedIn,
		"isBotOn":    isBotOn,

		"currentTwitchBroadcasterID": currentTwitchBroadcasterID,
		"currentTwitchBotUserID":     currentTwitchBotUserID,
		"currentTwitchBotUserName":   currentTwitchBotUserName,
	})
}
