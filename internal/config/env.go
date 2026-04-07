package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const EnvironmentPrefixMsg = "Environment variable "
const EnvironmentSuffixMsg = " is required."

type Env struct {
	GinMode                 string
	AppName                 string
	ServerPort              string
	TwitchClientName        string
	TwitchClientID          string
	TwitchClientSecret      string
	TwitchClientRedirectURI string
	TwitchBroadcasterID     string
	TwitchBotUserID         string
	TwitchBotUserName       string
}

var env *Env

func mustExistGoEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	if value != "debug" && value != "release" {
		log.Fatalf("%s must be debug or release", key)
	}

	return value
}

func mustExistString(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	return value
}

func InitEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	env = &Env{
		GinMode:                 mustExistGoEnv("GIN_MODE"),
		AppName:                 mustExistString("APP_NAME"),
		ServerPort:              mustExistString("SERVER_PORT"),
		TwitchClientName:        mustExistString("TWITCH_CLIENT_NAME"),
		TwitchClientID:          mustExistString("TWITCH_CLIENT_ID"),
		TwitchClientSecret:      mustExistString("TWITCH_CLIENT_SECRET"),
		TwitchClientRedirectURI: mustExistString("TWITCH_CLIENT_REDIRECT_URI"),
		TwitchBroadcasterID:     mustExistString("TWITCH_BROADCASTER_ID"),
		TwitchBotUserID:         mustExistString("TWITCH_BOT_USER_ID"),
		TwitchBotUserName:       mustExistString("TWITCH_BOT_USER_NAME"),
	}

	log.Println("Environment variables initialized successfully")

	return env
}

func GetEnv() *Env {
	if env == nil {
		log.Fatal("env not initialized, call InitEnv first")
	}

	return env
}

/**
func mustExistBool(key string) bool {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	return value == "true"
}

func mustExistInt(key string) int {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf(EnvironmentPrefixMsg+key+" must be a valid integer: %v", err)
	}

	return intValue
}

func mustExistStringSlice(key string) []string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	return strings.Split(value, ",")
}

func mustExistDuration(key string) time.Duration {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf(EnvironmentPrefixMsg+key+" must be a valid duration: %v", err)
	}

	return duration
}
*/
