package main

import (
	"encoding/json"
	"fmt"
	"os"

	chatgptgo "github.com/acheong08/ChatGPT-Go/chatbot"
)

var accessToken = os.Getenv("ACCESS_TOKEN")

func main() {
	chatbot, err := chatgptgo.NewChatbot(accessToken)
	if err != nil {
		panic(err)
	}
	history, err := chatbot.GetHistory(0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(history)

	conversationID := history.Items[0].ID
	conversation, err := chatbot.GetConversation(conversationID)
	if err != nil {
		panic(err)
	}
	jsonString, err := json.MarshalIndent(conversation, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonString))
}
