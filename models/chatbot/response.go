package chatbot

import "github.com/acheong08/ChatGPT-Go/models"

type ChatbotResponse struct {
	Message models.Message `json:"message"`
	ConversationID string `json:"conversation_id"`
	Error *string `json:"error"`
}
