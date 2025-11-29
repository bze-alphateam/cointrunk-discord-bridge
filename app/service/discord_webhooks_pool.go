package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bze-alphateam/cointrunk-discord-bridge/app/dto"
)

type DiscordWebhooksPool struct {
	webhooks []string
}

func NewDiscordWebhooksPool(webhooks []string) (*DiscordWebhooksPool, error) {
	if len(webhooks) == 0 {
		return nil, fmt.Errorf("invalid webhook url provided to DiscordWebhooksPool constructor")
	}

	return &DiscordWebhooksPool{webhooks: webhooks}, nil
}

func (w DiscordWebhooksPool) PostMessage(message dto.WebhookMessage) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for _, webhook := range w.webhooks {
		if err = w.postMessageOnWebhook(messageBytes, webhook); err != nil {
			log.Println(fmt.Sprintf("failed to post message to webhook %s: %s", webhook[:60], err))
		}
	}

	return nil
}

func (w DiscordWebhooksPool) postMessageOnWebhook(messageBytes []byte, webhook string) error {
	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(messageBytes))
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
