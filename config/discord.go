package config

import (
	"fmt"
	"os"
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
	if len(d.Webhook) == 0 {
		return fmt.Errorf("invalid discord webhook")
	}

	return nil
}
