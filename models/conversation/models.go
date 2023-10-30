package conversation

import "github.com/acheong08/ChatGPT-Go/models"

type Conversations struct {
	Items                   []ConversationBrief `json:"items"`
	Total                   int                 `json:"total"`
	Limit                   int                 `json:"limit"`
	Offset                  int                 `json:"offset"`
	HasConversationsMissing bool                `json:"has_conversations_missing"`
}

type ConversationBrief struct {
	ID                     string                  `json:"id"`
	Title                  string                  `json:"title"`
	CreateTime             string                  `json:"create_time"`
	UpdateTime             string                  `json:"update_time"`
	Mapping                *map[string]models.Node `json:"mapping"`
	CurrentNode            *string                 `json:"current_node"`
	ConversationTemplateID *string                 `json:"conversation_template_id"`
	ModerationResults      *[]any                  `json:"moderation_results"`
}

type Conversation struct {
	ConversationBrief
	CreateTime float32 `json:"create_time"`
	UpdateTime float32 `json:"update_time"`
}
