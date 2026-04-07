package main

import (
	"log"

	"github.com/gabriel-logan/twitch-ai-bot/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	router.RegisterRouter(r)

	serverPort := "8080"

	log.Println("Server running at: http://localhost:" + serverPort)
	if err := r.Run(":" + serverPort); err != nil {
		log.Printf("Error when trying to start server: %v \n", err)
	}
}
