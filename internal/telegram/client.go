package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	token  string
	apiURL string
	http   *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:  token,
		apiURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		http:   &http.Client{},
	}
}

func (c *Client) do(method string, req any) ([]byte, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.apiURL, method), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
			return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) SendMessage(chatID int64, text string, keyboard *InlineKeyboard) (int, error) {
	req := SendMessageReq{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: keyboard,
	}

	data, err := c.do("sendMessage", req)
	if err != nil {
		return 0, err
	}

	var apiResp APIResponse[Message]
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return 0, err
	}

	if !apiResp.OK {
		return 0, fmt.Errorf("sendMessage failed: %s", apiResp.Description)
	}

	return apiResp.Result.MessageID, nil
}

func (c *Client) EditMessage(chatID int64, msgID int, text string) error {
	req := EditMessageReq{
		ChatID:    chatID,
		MessageID: msgID,
		Text:      text,
	}

	data, err := c.do("editMessageText", req)
	if err != nil {
		return err
	}

	var apiResp APIResponse[Message]
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return err
	}

	if !apiResp.OK {
		return fmt.Errorf("editMessage failed: %s", apiResp.Description)
	}

	return nil
}

func (c *Client) AnswerCallback(callbackID string, text string) error {
	req := AnswerCallbackReq{
		CallbackQueryID: callbackID,
		Text:            text,
	}

	data, err := c.do("answerCallbackQuery", req)
	if err != nil {
		return err
	}

	var apiResp APIResponse[bool]
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return err
	}

	if !apiResp.OK {
		return fmt.Errorf("answerCallback failed: %s", apiResp.Description)
	}

	return nil
}

func (c *Client) GetUpdates(offset int, timeoutSec int) ([]Update, error) {
	req := GetUpdatesReq{
		Offset:  offset,
		Timeout: timeoutSec,
		AllowedUpdates: []string{"message", "callback_query"},
	}

	data, err := c.do("getUpdates", req)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse[[]Update]
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, err
	}

	if !apiResp.OK {
		return nil, fmt.Errorf("getUpdates failed: %s", apiResp.Description)
	}

	return apiResp.Result, nil
}