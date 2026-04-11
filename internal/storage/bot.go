package storage

var BotIsOn bool

func SetBotIsOn(isOn bool) {
	BotIsOn = isOn
}

func GetBotIsOn() bool {
	return BotIsOn
}
