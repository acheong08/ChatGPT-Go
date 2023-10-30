package chatgptgo

import (
	"github.com/bogdanfinn/tls-client"
)

type Chatbot struct {
	AccessToken string
	HTTPClient  *tls_client.HttpClient
}

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

type ChatbotRequest struct {
	Action           action `json:"action"`
	ArkoseToken      string `json:"arkose_token"`
	ConversationID   string `json:"conversation_id,omitempty"`
	ConversationMode struct {
		Kind string `json:"kind"`
	} `json:"conversation_mode"`
	ForceParagen               bool      `json:"force_paragen"`
	ForceRateLimit             bool      `json:"force_rate_limit"`
	HistoryAndTrainingDisabled bool      `json:"history_and_training_disabled"`
	Messages                   []message `json:"messages"`
	Model                      model     `json:"model"`
	ParentMessageID            string    `json:"parent_message_id"`
	Suggestions                []string  `json:"suggestions"`
	TimeZoneOffsetMin          int       `json:"time_zone_offset_min"`
}

type Conversations struct {
	Items                   []ConversationBrief `json:"items"`
	Total                   int                 `json:"total"`
	Limit                   int                 `json:"limit"`
	Offset                  int                 `json:"offset"`
	HasConversationsMissing bool                `json:"has_conversations_missing"`
}

type ConversationBrief struct {
	ID                     string           `json:"id"`
	Title                  string           `json:"title"`
	CreateTime             string           `json:"create_time"`
	UpdateTime             string           `json:"update_time"`
	Mapping                *map[string]node `json:"mapping"`
	CurrentNode            *string          `json:"current_node"`
	ConversationTemplateID *string          `json:"conversation_template_id"`
	ModerationResults      *[]any           `json:"moderation_results"`
}

type Conversation struct {
	ConversationBrief
	CreateTime float32 `json:"create_time"`
	UpdateTime float32 `json:"update_time"`
}

type node struct {
	Id       string    `json:"id"`
	Message  *message  `json:"message"`
	Parent   *string   `json:"parent"`
	Children *[]string `json:"children"`
}

type message struct {
	ID     string `json:"id"`
	Author struct {
		Role string `json:"role"`
	} `json:"author"`
	Content struct {
		ContentType string   `json:"content_type"`
		Parts       []string `json:"parts"`
	} `json:"content"`
	Metadata  map[string]any `json:"metadata"`
	Recipient string         `json:"recipient"`
}
