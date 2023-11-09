package chatgptgo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/acheong08/ChatGPT-Go/config"
	"github.com/acheong08/ChatGPT-Go/models/chatbot"
	"github.com/acheong08/ChatGPT-Go/models/conversation"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func NewChatbot(accessToken string) (Chatbot, error) {
	client, err := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithClientProfile(profiles.Firefox_117),
		tls_client.WithRandomTLSExtensionOrder(),
	)
	if config.PUID != "" {
		client.SetCookies(&url.URL{
			Scheme: "https",
			Host:   "chat.openai.com",
		}, []*http.Cookie{
			{
				Name:  "_puid",
				Value: config.PUID,
			},
		})
	}
	return Chatbot{
		AccessToken: accessToken,
		HTTPClient:  &client,
	}, err
}

func (c *Chatbot) GetHistory(offset, limit int) (*conversation.Conversations, error) {
	if limit == 0 {
		limit = 28
	}
	var conversations conversation.Conversations
	err := c.makeRequest(
		"GET",
		fmt.Sprintf("https://chat.openai.com/backend-api/conversations?offset=%d&limit=%d&order=updated", offset, limit),
		nil,
		&conversations,
	)
	return &conversations, err
}

func (c *Chatbot) GetConversation(conversationID string) (*conversation.Conversation, error) {
	var conversation conversation.Conversation
	err := c.makeRequest(
		"GET",
		fmt.Sprintf("https://chat.openai.com/backend-api/conversation/%s", conversationID),
		nil,
		&conversation,
	)
	return &conversation, err
}

// Ask No Stream
func (c *Chatbot) AskNS(body *chatbot.ChatbotRequest) (map[string]*chatbot.ChatbotResponse, error) {
	ch := make(chan chatbot.ChatbotResponse)
	cherr := make(chan error)
	go c.Ask(body, ch, cherr)
	resps := make(map[string]*chatbot.ChatbotResponse)
	for {
		select {
		case data := <-ch:
			tool_finished := false
			if status, ok := data.Message.Metadata["status"]; ok && status == "finished" {
				tool_finished = true
			}
			if data.Message.EndTurn || tool_finished || data.Message.Status == "finished_successfully" {
				role := data.Message.Author.Role
				if data.Message.Author.Name != "" {
					log.Println(data)
					role = data.Message.Author.Name
					if role == "dalle.text2im" {
						if data.Message.Content.ContentType != "multimodal_text" {
							continue
						}
					}
				}
				resps[role] = &data
			}
		case err := <-cherr:
			if err != nil {
				if err.Error() == config.ErrStreamEnd {
					return resps, nil
				} else {
					return nil, err
				}
			}
		}
	}
}

func (c *Chatbot) Ask(body *chatbot.ChatbotRequest, ch chan chatbot.ChatbotResponse, errch chan error) {
	if body == nil {
		errch <- fmt.Errorf("Body cannot be nil")
		return
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		errch <- err
		return
	}
	fmt.Println(string(bodyBytes))
	body_reader := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodPost, "https://chat.openai.com/backend-api/conversation", body_reader)
	if err != nil {
		errch <- err
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	addHeaders(req)
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		errch <- err
		return
	}
	if resp.StatusCode != 200 {
		// Read body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			errch <- err
			return
		}
		fmt.Println(string(bodyBytes))
		errch <- fmt.Errorf("Received status code %d", resp.StatusCode)
	}
	go func(body io.ReadCloser, ch chan chatbot.ChatbotResponse, errch chan error) {
		defer body.Close()
		defer close(ch)
		defer close(errch)
		// Stream response line by line
		scanner := bufio.NewScanner(body)
		for scanner.Scan() {
			var resp chatbot.ChatbotResponse
			raw := scanner.Text()
			if raw == "" {
				continue
			}
			if raw == "data: [DONE]" {
				errch <- fmt.Errorf(config.ErrStreamEnd)
				return
			}
			raw = strings.Replace(raw, "data: ", "", 1)
			err := json.Unmarshal([]byte(raw), &resp)
			if err != nil {
				fmt.Println(raw)
				fmt.Println(err)
				continue
			}
			ch <- resp

		}
		if err := scanner.Err(); err != nil {
			// Send error to error channel
			errch <- err
			return
		}
		// Send error to end stream
		errch <- fmt.Errorf(config.ErrStreamEnd)

	}(resp.Body, ch, errch)

}

func (c *Chatbot) GetModels() ([]chatbot.Model, error) {
 	var respModel struct {
		Models	[]chatbot.Model `json:"models"`
	}
	err := c.makeRequest(
		"GET",
		"https://chat.openai.com/backend-api/models",
		nil,
		&respModel,
	)
	return respModel.Models, err
}

func (c *Chatbot) makeRequest(method, url string, body, obj any) error {
	if method == "GET" && body != nil {
		return fmt.Errorf("Cannot send body with GET request")
	}
	var req *http.Request
	var err error
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		body_reader := bytes.NewReader(bodyBytes)
		req, err = http.NewRequest(method, url, body_reader)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	addHeaders(req)
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Received status code %d", resp.StatusCode)
	}
	if obj != nil {
		err = json.NewDecoder(resp.Body).Decode(obj)
		return err
	}
	return nil
}

type Download struct {
	Status      string `json:"status"`
	DownloadURL string `json:"download_url"`
	FileName    string `json:"file_name"`
}

func (c *Chatbot) DownloadFile(fileID string) (Download, error) {
	fileID = strings.Replace(fileID, "file-service://", "", 1)
	var download Download
	err := c.makeRequest(
		"GET",
		fmt.Sprintf("https://chat.openai.com/backend-api/files/%s/download", fileID),
		nil,
		&download,
	)
	return download, err
}

func addHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/118.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://chat.openai.com/")
	req.Header.Set("DNT", "1")
	req.Header.Set("Alt-Used", "chat.openai.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("TE", "trailers")
}
