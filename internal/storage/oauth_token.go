package storage

var OauthToken string

func SetOauthToken(token string) {
	OauthToken = token
}

func GetOauthToken() string {
	return OauthToken
}

func ClearOauthToken() {
	OauthToken = ""
}
