package ws

import "strings"

func hasMentionToBot(fragments []WSMessagePayloadMessageFragment, botUsername string) bool {
	for _, f := range fragments {
		if f.Mention != nil {
			if strings.EqualFold(f.Mention.UserLogin, botUsername) {
				return true
			}
		}
	}

	return false
}
