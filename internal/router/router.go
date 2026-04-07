package router

import (
	"github.com/gabriel-logan/twitch-ai-bot/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.Engine) {
	router.GET("/", handler.Index)
}
