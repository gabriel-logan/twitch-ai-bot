package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gabriel-logan/twitch-ai-bot/internal/ws"
	"github.com/gin-gonic/gin"
)

type CallbackTwitchResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func SignInTwitch(ctx *gin.Context) {
	env := config.GetEnv()

	const baseURL = "https://id.twitch.tv/oauth2/authorize"

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		log.Printf("Error when trying to create request: %v \n", err)
		ctx.JSON(http.StatusInternalServerError, "Error when trying to create request")
		return
	}

	q := req.URL.Query()
	q.Add("client_id", env.TwitchClientID)
	q.Add("redirect_uri", env.TwitchClientRedirectURI)
	q.Add("response_type", "code")
	q.Add("scope", "user:read:chat user:write:chat user:bot channel:bot")

	req.URL.RawQuery = q.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, req.URL.String())
}

func SignOutTwitch(ctx *gin.Context) {
	storage.ClearOauthToken()
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func CallbackTwitch(ctx *gin.Context) {
	code := ctx.Query("code")

	if code == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	const baseURL = "https://id.twitch.tv/oauth2/token"

	req, err := http.NewRequest(http.MethodPost, baseURL, nil)
	if err != nil {
		log.Printf("Error when trying to create request: %v \n", err)
		ctx.JSON(http.StatusInternalServerError, "Error when trying to create request")
		return
	}

	q := req.URL.Query()
	q.Add("client_id", config.GetEnv().TwitchClientID)
	q.Add("client_secret", config.GetEnv().TwitchClientSecret)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("redirect_uri", config.GetEnv().TwitchClientRedirectURI)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error when trying to create request: %v \n", err)
		ctx.JSON(http.StatusInternalServerError, "Error when trying to create request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error when trying to get access token: %v \n", resp)
		ctx.JSON(http.StatusInternalServerError, "Error when trying to get access token")
		return
	}

	var callbackTwitchResponse CallbackTwitchResponse
	if err := json.NewDecoder(resp.Body).Decode(&callbackTwitchResponse); err != nil {
		log.Printf("Error when trying to decode response: %v \n", err)
		ctx.JSON(http.StatusInternalServerError, "Error when trying to decode response")
		return
	}

	storage.SetOauthToken(callbackTwitchResponse.AccessToken)

	go ws.StartBot()

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}
