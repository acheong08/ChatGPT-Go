package chatbot

type Model struct {
	Slug string `json:"slug"`
	MaxTokens int `json:"max_tokens"`
	Title string `json:"title"`
	Description string `json:"description"`
	Tags []string `json:"tags"`
	EnabledTools []string `json:"enabled_tools"`
	ProductFeatures map[string]any `json:"product_features"`
}
