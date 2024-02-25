package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/dto"
	"net/http"
)

type DiscordWebhook struct {
	webhook string
}

func NewDiscordWebhook(webhook string) (*DiscordWebhook, error) {
	if len(webhook) == 0 {
		return nil, fmt.Errorf("invalid webhook url provided to DiscordWebhook constructor")
	}

	return &DiscordWebhook{webhook: webhook}, nil
}

func (w DiscordWebhook) PostMessage(message dto.WebhookMessage) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", w.webhook, bytes.NewBuffer(messageBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response received from discord webhook: %d", resp.StatusCode)
	}

	return nil
}
