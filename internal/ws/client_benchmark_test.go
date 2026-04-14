package ws

import (
	"testing"
	"time"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
)

func setupEnvValue() config.Env {
	return config.Env{
		GinMode:                        "debug",
		AppName:                        "twitch-ai-bot",
		ServerPort:                     "8080",
		ServerTrustedProxies:           []string{"127.0.0.1"},
		TwitchClientID:                 "test-client-id",
		TwitchClientSecret:             "test-client-secret",
		TwitchClientRedirectURI:        "http://localhost:8080/callback",
		TwitchBroadcasterID:            "123456789",
		TwitchBotUserID:                "987654321",
		TwitchBotUserName:              "testbot",
		TwitchKeyWordToCallBot:         "ai",
		TwitchChatMessageMaxLength:     480,
		TwitchTimeforTheBottoTellaJoke: 5 * time.Minute,
		GroqAPIKey:                     "test-groq-api-key",
		GroqModels:                     []string{"mixtral-8x7b-32768", "llama-3.3-70b-versatile"},
		GroqMaxContextInput:            10,
		ContextRequestDuration:         30 * time.Second,
	}
}

func setupEnvPtr() *config.Env {
	return &config.Env{
		GinMode:                        "debug",
		AppName:                        "twitch-ai-bot",
		ServerPort:                     "8080",
		ServerTrustedProxies:           []string{"127.0.0.1"},
		TwitchClientID:                 "test-client-id",
		TwitchClientSecret:             "test-client-secret",
		TwitchClientRedirectURI:        "http://localhost:8080/callback",
		TwitchBroadcasterID:            "123456789",
		TwitchBotUserID:                "987654321",
		TwitchBotUserName:              "testbot",
		TwitchKeyWordToCallBot:         "ai",
		TwitchChatMessageMaxLength:     480,
		TwitchTimeforTheBottoTellaJoke: 5 * time.Minute,
		GroqAPIKey:                     "test-groq-api-key",
		GroqModels:                     []string{"mixtral-8x7b-32768", "llama-3.3-70b-versatile"},
		GroqMaxContextInput:            10,
		ContextRequestDuration:         30 * time.Second,
	}
}

var (
	sinkString   string
	sinkInt      int
	sinkDuration time.Duration
	sinkSliceLen int
)

func consumeEnvValue(env config.Env) {
	sinkString = env.GinMode
	sinkString = env.AppName
	sinkString = env.ServerPort
	sinkSliceLen = len(env.ServerTrustedProxies)
	sinkString = env.TwitchClientID
	sinkString = env.TwitchClientSecret
	sinkString = env.TwitchClientRedirectURI
	sinkString = env.TwitchBroadcasterID
	sinkString = env.TwitchBotUserID
	sinkString = env.TwitchBotUserName
	sinkString = env.TwitchKeyWordToCallBot
	sinkInt = env.TwitchChatMessageMaxLength
	sinkDuration = env.TwitchTimeforTheBottoTellaJoke
	sinkString = env.GroqAPIKey
	sinkSliceLen = len(env.GroqModels)
	sinkInt = env.GroqMaxContextInput
	sinkDuration = env.ContextRequestDuration
}

func consumeEnvPtr(env *config.Env) {
	sinkString = env.GinMode
	sinkString = env.AppName
	sinkString = env.ServerPort
	sinkSliceLen = len(env.ServerTrustedProxies)
	sinkString = env.TwitchClientID
	sinkString = env.TwitchClientSecret
	sinkString = env.TwitchClientRedirectURI
	sinkString = env.TwitchBroadcasterID
	sinkString = env.TwitchBotUserID
	sinkString = env.TwitchBotUserName
	sinkString = env.TwitchKeyWordToCallBot
	sinkInt = env.TwitchChatMessageMaxLength
	sinkDuration = env.TwitchTimeforTheBottoTellaJoke
	sinkString = env.GroqAPIKey
	sinkSliceLen = len(env.GroqModels)
	sinkInt = env.GroqMaxContextInput
	sinkDuration = env.ContextRequestDuration
}

func BenchmarkSetupEnvValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		env := setupEnvValue()

		sinkString = env.GinMode
		sinkString = env.AppName
		sinkString = env.ServerPort
		sinkSliceLen = len(env.ServerTrustedProxies)
		sinkString = env.TwitchClientID
		sinkString = env.TwitchClientSecret
		sinkString = env.TwitchClientRedirectURI
		sinkString = env.TwitchBroadcasterID
		sinkString = env.TwitchBotUserID
		sinkString = env.TwitchBotUserName
		sinkString = env.TwitchKeyWordToCallBot
		sinkInt = env.TwitchChatMessageMaxLength
		sinkDuration = env.TwitchTimeforTheBottoTellaJoke
		sinkString = env.GroqAPIKey
		sinkSliceLen = len(env.GroqModels)
		sinkInt = env.GroqMaxContextInput
		sinkDuration = env.ContextRequestDuration
	}
}

func BenchmarkSetupEnvPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		env := setupEnvPtr()

		sinkString = env.GinMode
		sinkString = env.AppName
		sinkString = env.ServerPort
		sinkSliceLen = len(env.ServerTrustedProxies)
		sinkString = env.TwitchClientID
		sinkString = env.TwitchClientSecret
		sinkString = env.TwitchClientRedirectURI
		sinkString = env.TwitchBroadcasterID
		sinkString = env.TwitchBotUserID
		sinkString = env.TwitchBotUserName
		sinkString = env.TwitchKeyWordToCallBot
		sinkInt = env.TwitchChatMessageMaxLength
		sinkDuration = env.TwitchTimeforTheBottoTellaJoke
		sinkString = env.GroqAPIKey
		sinkSliceLen = len(env.GroqModels)
		sinkInt = env.GroqMaxContextInput
		sinkDuration = env.ContextRequestDuration
	}
}

func BenchmarkPassEnvValue(b *testing.B) {
	env := setupEnvValue()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		consumeEnvValue(env)
	}
}

func BenchmarkPassEnvPtr(b *testing.B) {
	env := setupEnvPtr()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		consumeEnvPtr(env)
	}
}

func BenchmarkAccessAllFieldsValue(b *testing.B) {
	env := setupEnvValue()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sinkString = env.GinMode
		sinkString = env.AppName
		sinkString = env.ServerPort
		sinkSliceLen = len(env.ServerTrustedProxies)
		sinkString = env.TwitchClientID
		sinkString = env.TwitchClientSecret
		sinkString = env.TwitchClientRedirectURI
		sinkString = env.TwitchBroadcasterID
		sinkString = env.TwitchBotUserID
		sinkString = env.TwitchBotUserName
		sinkString = env.TwitchKeyWordToCallBot
		sinkInt = env.TwitchChatMessageMaxLength
		sinkDuration = env.TwitchTimeforTheBottoTellaJoke
		sinkString = env.GroqAPIKey
		sinkSliceLen = len(env.GroqModels)
		sinkInt = env.GroqMaxContextInput
		sinkDuration = env.ContextRequestDuration
	}
}

func BenchmarkAccessAllFieldsPtr(b *testing.B) {
	env := setupEnvPtr()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sinkString = env.GinMode
		sinkString = env.AppName
		sinkString = env.ServerPort
		sinkSliceLen = len(env.ServerTrustedProxies)
		sinkString = env.TwitchClientID
		sinkString = env.TwitchClientSecret
		sinkString = env.TwitchClientRedirectURI
		sinkString = env.TwitchBroadcasterID
		sinkString = env.TwitchBotUserID
		sinkString = env.TwitchBotUserName
		sinkString = env.TwitchKeyWordToCallBot
		sinkInt = env.TwitchChatMessageMaxLength
		sinkDuration = env.TwitchTimeforTheBottoTellaJoke
		sinkString = env.GroqAPIKey
		sinkSliceLen = len(env.GroqModels)
		sinkInt = env.GroqMaxContextInput
		sinkDuration = env.ContextRequestDuration
	}
}

func BenchmarkPassEnvValueThroughLayers(b *testing.B) {
	env := setupEnvValue()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		func(env config.Env) {
			func(env config.Env) {
				func(env config.Env) {
					func(env config.Env) {
						consumeEnvValue(env)
					}(env)
				}(env)
			}(env)
		}(env)
	}
}

func BenchmarkPassEnvPtrThroughLayers(b *testing.B) {
	env := setupEnvPtr()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		func(env *config.Env) {
			func(env *config.Env) {
				func(env *config.Env) {
					func(env *config.Env) {
						consumeEnvPtr(env)
					}(env)
				}(env)
			}(env)
		}(env)
	}
}
