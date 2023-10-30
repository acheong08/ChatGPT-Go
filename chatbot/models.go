package chatbot

import (
	"github.com/bogdanfinn/tls-client"
)

type Chatbot struct {
	AccessToken string
	HTTPClient  *tls_client.HttpClient
}

