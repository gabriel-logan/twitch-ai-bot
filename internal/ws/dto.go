package ws

type WSMessageMetadata struct {
	MessageType      string `json:"message_type"`
	SubscriptionType string `json:"subscription_type"`
}

type WSMessagePayloadMessageFragment struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Mention *struct {
		UserID    string `json:"user_id"`
		UserLogin string `json:"user_login"`
		UserName  string `json:"user_name"`
	} `json:"mention"`
}

type WSMessagePayloadMessage struct {
	Text      string                            `json:"text"`
	Fragments []WSMessagePayloadMessageFragment `json:"fragments"`
}

type WSMessagePayload struct {
	Session struct {
		ID string `json:"id"`
	} `json:"session"`
	Event struct {
		BroadcasterUserLogin string                  `json:"broadcaster_user_login"`
		ChatterUserLogin     string                  `json:"chatter_user_login"`
		SystemMessage        string                  `json:"system_message"`
		Message              WSMessagePayloadMessage `json:"message"`
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
