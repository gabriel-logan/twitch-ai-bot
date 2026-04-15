package config

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const EnvironmentPrefixMsg = "Environment variable "
const EnvironmentSuffixMsg = " is required."

type Env struct {
	GinMode                        string
	AppName                        string
	ServerPort                     string
	ServerTrustedProxies           []string
	TwitchClientID                 string
	TwitchClientSecret             string
	TwitchClientRedirectURI        string
	TwitchBroadcasterID            string
	TwitchBotUserID                string
	TwitchBotUserName              string
	TwitchKeyWordToCallBot         string
	TwitchChatMessageMaxLength     int
	TwitchTimeforTheBottoTellaJoke time.Duration
	GroqAPIKey                     string
	GroqModels                     []string
	GroqMaxContextInput            int
	ContextRequestDuration         time.Duration
}

var env *Env

func mustExistBool(key string) (bool, error) {
	value := os.Getenv(key)

	if value == "" {
		return false, errors.New(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	isBool, err := strconv.ParseBool(value)
	if err != nil {
		return false, errors.New(EnvironmentPrefixMsg + key + " must be a valid boolean: " + err.Error())
	}

	return isBool, nil
}

func mustExistGoEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", errors.New(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	if value != "debug" && value != "release" {
		return "", errors.New(EnvironmentPrefixMsg + key + " must be debug or release")
	}

	return value, nil
}

func mustExistString(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", errors.New(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	return value, nil
}

func mustExistInt(key string) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return 0, errors.New(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New(EnvironmentPrefixMsg + key + " must be a valid integer: " + err.Error())
	}

	return intValue, nil
}

func mustExistStringSlice(key string) ([]string, error) {
	value := os.Getenv(key)

	if value == "" {
		return nil, errors.New(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	return strings.Split(value, ","), nil
}

func mustExistDuration(key string) (time.Duration, error) {
	value := os.Getenv(key)

	if value == "" {
		return 0, errors.New(EnvironmentPrefixMsg + key + EnvironmentSuffixMsg)
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, errors.New(EnvironmentPrefixMsg + key + " must be a valid duration: " + err.Error())
	}

	return duration, nil
}

func loadEnv() *Env { // nosonar
	e := &Env{}

	t := reflect.TypeOf(e).Elem()
	envLength := t.NumField()

	var checksDone int
	var err error
	errs := make([]error, 0, envLength)

	if e.GinMode, err = mustExistGoEnv("GIN_MODE"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.AppName, err = mustExistString("APP_NAME"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.ServerPort, err = mustExistString("SERVER_PORT"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.ServerTrustedProxies, err = mustExistStringSlice("SERVER_TRUSTED_PROXIES"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchClientID, err = mustExistString("TWITCH_CLIENT_ID"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchClientSecret, err = mustExistString("TWITCH_CLIENT_SECRET"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchClientRedirectURI, err = mustExistString("TWITCH_CLIENT_REDIRECT_URI"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchBroadcasterID, err = mustExistString("TWITCH_BROADCASTER_ID"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchBotUserID, err = mustExistString("TWITCH_BOT_USER_ID"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchBotUserName, err = mustExistString("TWITCH_BOT_USER_NAME"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchKeyWordToCallBot, err = mustExistString("TWITCH_KEY_WORD_TO_CALL_BOT"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchChatMessageMaxLength, err = mustExistInt("TWITCH_CHAT_MESSAGE_MAX_LENGTH"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.TwitchTimeforTheBottoTellaJoke, err = mustExistDuration("TWITCH_TIME_FOR_THE_BOT_TO_TELLA_JOKE"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.GroqAPIKey, err = mustExistString("GROQ_API_KEY"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.GroqModels, err = mustExistStringSlice("GROQ_MODELS"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.GroqMaxContextInput, err = mustExistInt("GROQ_MAX_CONTEXT_INPUT"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if e.ContextRequestDuration, err = mustExistDuration("CONTEXT_REQUEST_DURATION"); err != nil {
		errs = append(errs, err)
	}
	checksDone++

	if checksDone != envLength {
		log.Printf("Warning: Number of environment variables checked (%d) does not match number of fields in Env struct (%d). This may indicate that some environment variables are not being checked.\n", checksDone, envLength)
	}

	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}

		log.Fatal("Failed to load environment variables")
	}

	return e
}

func InitEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	env = loadEnv()

	log.Println("Environment variables initialized successfully")

	return env
}

func ReloadEnv() {
	env = loadEnv()

	log.Println("Environment variables reloaded successfully")
}

func GetEnv() *Env {
	if env == nil {
		log.Fatal("env not initialized, call InitEnv first")
	}

	return env
}
