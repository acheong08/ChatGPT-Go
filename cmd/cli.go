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
	// history, err := cb.GetHistory(0, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(history)
	//
	// conversationID := history.Items[0].ID
	// conversation, err := cb.GetConversation(conversationID)
	// if err != nil {
	// 	panic(err)
	// }
	// jsonString, err := json.MarshalIndent(conversation, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(jsonString))
	// Get arkose token
	// Get models
	availableModels, err := cb.GetModels()
	if err != nil {
		panic(err)
	}
	// Check if "GPT-4 (All Tools)" is in the models
	var found bool
	for _, model := range availableModels {
		if model.Title == "GPT-4 (All Tools)" {
			found = true
			break
		}
	}
	if !found {
		panic("GPT-4 (All Tools) not found")
	}

	arkoseToken, err := captchaSolver.GetOpenAIToken(funcaptcha.ArkVerChat4, config.PUID)
	if err != nil {
		panic(err)
	}
	// Send a request
	req, err := chatbot.NewRequest(
		chatbot.WithMessage(
			models.NewMessage("user", models.MessageContent{
				ContentType: "text",
				Parts:       []any{"Generate an image of a creepy bird"},
			}),
		),
		chatbot.WithModel(chatbot.ModelGPT4),
		chatbot.WithArkoseToken(arkoseToken),
	)
	if err != nil {
		panic(err)
	}

	resps, err := cb.AskNS(req)
	if err != nil {
		panic(err)
	}

	for _, resp := range resps {
		contentType := resp.Message.Content.Type()
		fmt.Println(contentType)
		if contentType == "dalle_image" {
			fmt.Println("Got an image!")
			for _, part := range resp.Message.Content.Parts {
				fmt.Println(part)
				down, err := cb.DownloadFile(part.(models.DallEPart).AssetPointer)
				if err != nil {
					panic(err)
				}
				fmt.Println(down)
			}
		} else {
			fmt.Printf("Got a %s which is not dalle_image", contentType)
		}
	}

	jsonString, err := json.MarshalIndent(resps, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonString))

}
