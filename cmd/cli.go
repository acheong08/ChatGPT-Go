package main

import (
	"encoding/json"
	"fmt"
	"os"

	chatgptgo "github.com/acheong08/ChatGPT-Go"
	"github.com/acheong08/ChatGPT-Go/config"
	"github.com/acheong08/ChatGPT-Go/models"
	"github.com/acheong08/ChatGPT-Go/models/chatbot"
	"github.com/acheong08/funcaptcha"
)

var accessToken = os.Getenv("ACCESS_TOKEN")

func main() {
	captchaSolver := funcaptcha.NewSolver(funcaptcha.WithHarpool)
	cb, err := chatgptgo.NewChatbot(accessToken)
	if err != nil {
		panic(err)
	}
	history, err := cb.GetHistory(0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(history)

	conversationID := history.Items[0].ID
	conversation, err := cb.GetConversation(conversationID)
	if err != nil {
		panic(err)
	}
	jsonString, err := json.MarshalIndent(conversation, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonString))
	// Get arkose token
	arkoseToken, err := captchaSolver.GetOpenAIToken(funcaptcha.ArkVerChat4, config.PUID)
	if err != nil {
		panic(err)
	}
	// Send a request
	req, err := chatbot.NewRequest(
		chatbot.WithMessage(
			models.NewMessage("user", models.MessageContent{
				ContentType: "text",
				Parts:       []string{"Who won the most recent world cup?"},
			}),
		),
		chatbot.WithModel(chatbot.ModelBrowsing),
		chatbot.WithArkoseToken(arkoseToken),
	)
	if err != nil {
		panic(err)
	}
	ch := make(chan chatbot.ChatbotResponse)
	cherr := make(chan error)
	cb.Ask(req, ch, cherr)

	var stop bool
	for {
		select {
		case _ = <-ch:
			fmt.Print(".")
		case err := <-cherr:
			if err != nil {
				if err.Error() == config.ErrStreamEnd {
					stop = true
				} else {
					panic(err)
				}
			}
		}
		if stop {
			break
		}
	}

}
