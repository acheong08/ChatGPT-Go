package chatbot

import (
	"github.com/acheong08/ChatGPT-Go/models"
	"github.com/google/uuid"
)

type action string

const (
	ActionNext     action = "next"
	ActionContinue action = "continue"
)

type model string

const (
	ModelFree            model = "text-davinci-002-render-sha"
	ModelGPT4            model = "gpt-4"
	ModelBrowsing        model = "gpt-4-browsing"
	ModelCodeInterpreter model = "gpt-4-code-interpreter"
	ModelDalle           model = "gpt-4-dalle"
)

type conversationModes string

const (
	ConversationModePrimaryAssistant conversationModes = "primary_assistant"
)

type ChatbotRequest struct {
	Action                     action           `json:"action"`
	ArkoseToken                string           `json:"arkose_token"`
	ConversationID             string           `json:"conversation_id,omitempty"`
	ConversationMode           conversationMode `json:"conversation_mode"`
	ForceParagen               bool             `json:"force_paragen"`
	ForceRateLimit             bool             `json:"force_rate_limit"`
	HistoryAndTrainingDisabled bool             `json:"history_and_training_disabled"`
	Messages                   []models.Message `json:"messages"`
	Model                      model            `json:"model"`
	ParentMessageID            string           `json:"parent_message_id"`
	Suggestions                []string         `json:"suggestions"`
	TimeZoneOffsetMin          int              `json:"time_zone_offset_min"`
}

type conversationMode struct {
	Kind conversationModes `json:"kind"`
}

func defaultRequest() ChatbotRequest {
	return ChatbotRequest{
		Action:            ActionNext,
		ConversationMode:  conversationMode{Kind: ConversationModePrimaryAssistant},
		ForceParagen:      false,
		ForceRateLimit:    false,
		Messages:          []models.Message{},
		Model:             ModelFree,
		TimeZoneOffsetMin: 0,
		ParentMessageID:   uuid.New().String(),
	}
}

type requestArgs func(*ChatbotRequest)

func WithAction(action action) requestArgs {
	return func(request *ChatbotRequest) {
		request.Action = action
	}
}

func WithArkoseToken(arkoseToken string) requestArgs {
	return func(request *ChatbotRequest) {
		request.ArkoseToken = arkoseToken
	}
}

func WithConversationID(conversationID string) requestArgs {
	return func(request *ChatbotRequest) {
		request.ConversationID = conversationID
	}
}

func WithParentID(parentID string) requestArgs {
	return func(request *ChatbotRequest) {
		request.ParentMessageID = parentID
	}
}

func WithModel(model model) requestArgs {
	return func(request *ChatbotRequest) {
		request.Model = model
	}
}

func WithMessage(message models.Message) requestArgs {
	return func(request *ChatbotRequest) {
		request.Messages = append(request.Messages, message)
	}
}

func NewRequest(args ...requestArgs) (*ChatbotRequest, error) {
	request := defaultRequest()
	for _, arg := range args {
		arg(&request)
	}
	return &request, nil
}
