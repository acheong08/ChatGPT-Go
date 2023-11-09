package models

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/google/uuid"
)

type Node struct {
	Id       string    `json:"id"`
	Message  *Message  `json:"message"`
	Parent   *string   `json:"parent"`
	Children *[]string `json:"children"`
}

type Message struct {
	ID        string         `json:"id"`
	Author    Author         `json:"author"`
	Content   MessageContent `json:"content"`
	Metadata  map[string]any `json:"metadata"`
	Recipient string         `json:"recipient"`
	Status    string         `json:"status,omitempty"`
	EndTurn   bool           `json:"end_turn"`
}

type Author struct {
	Role string `json:"role"`
	Name string `json:"name,omitempty"`
}

type DallEPart struct {
	ContentType  string        `json:"content_type"`
	AssetPointer string        `json:"asset_pointer"`
	SizeBytes    float64       `json:"size_bytes"`
	Metadata     DallEMetadata `json:"metadata"`
}
type DallEMetadata struct {
	DallE struct {
		Prompt string `json:"prompt"`
	} `json:"dalle"`
}

type MessageContent struct {
	ContentType string `json:"content_type"`
	Parts       []any  `json:"parts"`
}

func (m *MessageContent) Type() string {
	if len(m.Parts) == 0 {
		return "empty"
	}
	typeOf := reflect.TypeOf(m.Parts[0])
	log.Println(fmt.Sprintf("TypeOf: %v", typeOf))
	if typeOf.Kind() == reflect.String {
		return "string"
	}
	// Check if it's a map
	if typeOf.Kind() == reflect.Map {
		// Check if it's a DALL-E image part
		if contentType, ok := m.Parts[0].(map[string]any)["content_type"]; ok && contentType == "image_asset_pointer" {
			// Convert from map to DallEPart
			for i := range m.Parts {
				jsonString, _ := json.Marshal(m.Parts[i])
				var part DallEPart
				json.Unmarshal(jsonString, &part)
				m.Parts[i] = part
			}

			return "dalle_image"
		}
	}
	return "unknown"
}
func NewMessage(role string, contents MessageContent) Message {
	// Generate message ID
	id := uuid.New().String()
	return Message{
		ID: id,
		Author: Author{
			Role: role,
		},
		Content:  contents,
		Metadata: map[string]any{},
	}
}
