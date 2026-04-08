package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gin-gonic/gin"
)

type CallbackTwitchResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func SignInTwitch(c *gin.Context) {
	env := config.GetEnv()

	ctx, cancel := context.WithTimeout(c.Request.Context(), env.ContextRequestDuration)
	defer cancel()

	const baseURL = "https://id.twitch.tv/oauth2/authorize"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		log.Printf("Error when trying to create request: %v \n", err)
		c.JSON(http.StatusInternalServerError, "Error when trying to create request")
		return
	}

	q := req.URL.Query()
	q.Add("client_id", env.TwitchClientID)
	q.Add("redirect_uri", env.TwitchClientRedirectURI)
	q.Add("response_type", "code")
	q.Add("scope", "user:read:chat user:write:chat user:bot channel:bot")

	req.URL.RawQuery = q.Encode()

	c.Redirect(http.StatusTemporaryRedirect, req.URL.String())
}

func SignOutTwitch(c *gin.Context) {
	storage.ClearOauthToken()

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func CallbackTwitch(c *gin.Context) {
	env := config.GetEnv()

	ctx, cancel := context.WithTimeout(c.Request.Context(), env.ContextRequestDuration)
	defer cancel()

	code := c.Query("code")

	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	const baseURL = "https://id.twitch.tv/oauth2/token"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, nil)
	if err != nil {
		log.Printf("Error when trying to create request: %v \n", err)
		c.JSON(http.StatusInternalServerError, "Error when trying to create request")
		return
	}

	q := req.URL.Query()
	q.Add("client_id", env.TwitchClientID)
	q.Add("client_secret", env.TwitchClientSecret)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("redirect_uri", env.TwitchClientRedirectURI)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error when trying to create request: %v \n", err)
		c.JSON(http.StatusInternalServerError, "Error when trying to create request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error when trying to get access token: %v \n", resp)
		c.JSON(http.StatusInternalServerError, "Error when trying to get access token")
		return
	}

	var callbackTwitchResponse CallbackTwitchResponse
	if err := json.NewDecoder(resp.Body).Decode(&callbackTwitchResponse); err != nil {
		log.Printf("Error when trying to decode response: %v \n", err)
		c.JSON(http.StatusInternalServerError, "Error when trying to decode response")
		return
	}

	storage.SetOauthToken(callbackTwitchResponse.AccessToken)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
