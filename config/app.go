package config

import "fmt"

type App struct {
	Discord Discord `yaml:"discord"`
	History History `yaml:"history"`
	RestUrl string  `yaml:"rest_url"`
	RpcUrl  string  `yaml:"rpc_url"`
}

func NewAppConfig() App {
	discord := NewDiscordConfig()
	history := NewHistoryConfig()
	app := App{
		Discord: discord,
		History: history,
	}

	return app
}

func (a App) Validate() error {
	if err := a.Discord.Validate(); err != nil {
		return err
	}

	if err := a.History.Validate(); err != nil {
		return err
	}

	if a.RestUrl == "" {
		return fmt.Errorf("invalid rest_url provided in config")
	}

	return nil
}
