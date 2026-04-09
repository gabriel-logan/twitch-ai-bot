package main

import (
	"log"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger - uses log package internally
	config.InitLogger()

	// Load environment variables
	env := config.InitEnv()
	gin.SetMode(env.GinMode)

	r := gin.Default()

	r.SetTrustedProxies(env.ServerTrustedProxies)

	r.LoadHTMLGlob("templates/*")

	r.Static("/public", "./public")

	router.RegisterRouter(r)

	log.Println("Server running at: http://localhost:" + env.ServerPort)
	log.Println("Application name: " + env.AppName)
	log.Printf("Running in %s mode \n", env.GinMode)
	if err := r.Run(":" + env.ServerPort); err != nil {
		log.Printf("Error when trying to start server: %v \n", err)
	}
}
