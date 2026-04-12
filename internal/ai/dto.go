package ai

type RequestMessage struct {
	Role    string `json:"role"`
	Name    string `json:"name,omitempty"` // Optional
	Content string `json:"content"`
}

type Request struct {
	Model    string           `json:"model"`
	Messages []RequestMessage `json:"messages"`
}

type ResponseChoiceMessage struct {
	Role    string  `json:"role"`
	Content string `json:"content"` // Optional - Groq returns string or null
}

type ResponseChoice struct {
	Message ResponseChoiceMessage `json:"message"`
}

type Response struct {
	Model   string           `json:"model"`
	Choices []ResponseChoice `json:"choices"`
}
