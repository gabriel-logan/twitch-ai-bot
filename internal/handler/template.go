package handler

import (
	"net/http"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	env := config.GetEnv()

	var IsLoggedIn bool
	if storage.OauthToken != "" {
		IsLoggedIn = true
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":      env.AppName,
		"isLoggedIn": IsLoggedIn,
	})
}
