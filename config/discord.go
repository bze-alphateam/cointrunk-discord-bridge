package config

import (
	"fmt"
	"os"
	"strings"
)

type Discord struct {
	Token            string `yaml:"token"`
	ChannelID        string `yaml:"channel_id"`
	Webhook          string `yaml:"webhook"`
	RoleChangerToken string `yaml:"role_changer_token"`
	AppID            string `yaml:"app_id"`
	GuildID          string `yaml:"guild_id"`
	RoleID           string `yaml:"role_id"`
}

func NewDiscordConfig() Discord {
	dToken := os.Getenv("DISCORD_TOKEN")
	dChannelId := os.Getenv("DISCORD_CHANNEL_ID")
	roleChangerToken := os.Getenv("DISCORD_ROLE_CHANGER_TOKEN")
	appId := os.Getenv("DISCORD_APP_ID")
	guildId := os.Getenv("DISCORD_GUILD_ID")
	roleId := os.Getenv("DISCORD_ROLE_ID")

	return Discord{
		Token:            dToken,
		ChannelID:        dChannelId,
		RoleChangerToken: roleChangerToken,
		AppID:            appId,
		GuildID:          guildId,
		RoleID:           roleId,
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

func (d Discord) ValidateRoleChanger() error {
	if d.RoleChangerToken == "" {
		return fmt.Errorf("invalid role_changer_token")
	}
	if d.AppID == "" {
		return fmt.Errorf("invalid app_id")
	}
	if d.GuildID == "" {
		return fmt.Errorf("invalid guild_id")
	}
	if d.RoleID == "" {
		return fmt.Errorf("invalid role_id")
	}
	return nil
}
