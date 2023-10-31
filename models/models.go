package models

import "github.com/google/uuid"

type Node struct {
	Id       string    `json:"id"`
	Message  *Message  `json:"message"`
	Parent   *string   `json:"parent"`
	Children *[]string `json:"children"`
}

type Message struct {
	ID        string           `json:"id"`
	Author    Author           `json:"author"`
	Content   MessageContent `json:"content"`
	Metadata  map[string]any   `json:"metadata"`
	Recipient string           `json:"recipient"`
}

type Author struct {
	Role string `json:"role"`
}

type MessageContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

func NewMessage(role string, contents MessageContent) Message {
	// Generate message ID
	id := uuid.New().String()
	return Message{
		ID: id,
		Author: Author{
			Role: role,
		},
		Content: contents,
		Metadata: map[string]any{},
	}
}
