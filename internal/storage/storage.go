package storage

var (
	OauthToken string
	BotIsOn    bool
)

func SetOauthToken(token string) {
	OauthToken = token
}

func GetOauthToken() string {
	return OauthToken
}

func ClearOauthToken() {
	OauthToken = ""
}

func SetBotIsOn(isOn bool) {
	BotIsOn = isOn
}

func GetBotIsOn() bool {
	return BotIsOn
}
