package ws

type WSMessageMetadata struct {
	MessageType      string `json:"message_type"`
	SubscriptionType string `json:"subscription_type"`
}

type WSMessagePayload struct {
	Session struct {
		ID string `json:"id"`
	} `json:"session"`
	Event struct {
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		ChatterUserLogin     string `json:"chatter_user_login"`
		Message              struct {
			Text string `json:"text"`
		} `json:"message"`
	} `json:"event"`
}

type WSMessage struct {
	Metadata WSMessageMetadata `json:"metadata"`
	Payload  WSMessagePayload  `json:"payload"`
}

type EventSubRequestCondition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
	UserID            string `json:"user_id"`
}

type EventSubRequestTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id"`
}

type EventSubRequest struct {
	Type      string                   `json:"type"`
	Version   string                   `json:"version"`
	Condition EventSubRequestCondition `json:"condition"`
	Transport EventSubRequestTransport `json:"transport"`
}

type SendMessageRequest struct {
	BroadcasterID string `json:"broadcaster_id"`
	SenderID      string `json:"sender_id"`
	Message       string `json:"message"`
}
