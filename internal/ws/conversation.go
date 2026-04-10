package ws

import "github.com/gabriel-logan/twitch-ai-bot/internal/ai"

type Conversation struct {
	System   ai.RequestMessage
	Messages []ai.RequestMessage
	Max      int
	Start    int
	Size     int
}

func NewConversation(max int, system ai.RequestMessage) *Conversation {
	return &Conversation{
		System:   system,
		Messages: make([]ai.RequestMessage, max),
		Max:      max,
		Start:    0,
		Size:     0,
	}
}

func (c *Conversation) Add(message ai.RequestMessage) {
	if message.Role == "system" {
		return
	}

	idx := (c.Start + c.Size) % c.Max

	c.Messages[idx] = message

	if c.Size < c.Max {
		c.Size++
	} else {
		c.Start = (c.Start + 1) % c.Max
	}
}

func (c *Conversation) Clear() {
	c.Messages = make([]ai.RequestMessage, c.Max)
	c.Start = 0
	c.Size = 0
}

func (c *Conversation) GetAll() []ai.RequestMessage {
	result := make([]ai.RequestMessage, c.Size)

	for i := 0; i < c.Size; i++ {
		result[i] = c.Messages[(c.Start+i)%c.Max]
	}

	return result
}

func (c *Conversation) BuildMessages() []ai.RequestMessage {
	messages := make([]ai.RequestMessage, 0, c.Size+1)

	messages = append(messages, c.System)
	messages = append(messages, c.GetAll()...)

	return messages
}
