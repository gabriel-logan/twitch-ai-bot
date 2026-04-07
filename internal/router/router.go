package router

import (
	"github.com/gabriel-logan/twitch-ai-bot/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.Engine) {
	router.GET("/", handler.Index)

	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/auth/sign-in/twitch", handler.SignInTwitch)
		apiRouter.GET("/auth/callback/twitch", handler.CallbackTwitch)
		apiRouter.GET("/auth/sign-out/twitch", handler.SignOutTwitch)

		apiRouter.GET("/twitch/user-info", handler.GetTwitchUserInfo)
	}
}
