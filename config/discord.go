package config

import (
	"fmt"
	"os"
	"strings"
)

type Discord struct {
	Token     string `yaml:"token"`
	ChannelID string `yaml:"channel_id"`
	Webhook   string `yaml:"webhook"`
}

func NewDiscordConfig() Discord {
	dToken := os.Getenv("DISCORD_TOKEN")
	dChannelId := os.Getenv("DISCORD_CHANNEL_ID")

	return Discord{
		Token:     dToken,
		ChannelID: dChannelId,
	}
}

func (d Discord) Validate() error {
	if len(d.GetWebhooks()) == 0 {
		return fmt.Errorf("invalid discord webhook")
	}

	return nil
}

func (d Discord) GetWebhooks() []string {
	split := strings.Split(d.Webhook, ",")
	var webhooks []string
	for _, s := range split {
		normalized := strings.TrimSpace(s)
		if len(normalized) < 60 {
			continue
		}

		webhooks = append(webhooks, normalized)
	}

	return webhooks
}
