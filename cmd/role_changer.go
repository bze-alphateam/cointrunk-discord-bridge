package cmd

import (
	"log"

	"github.com/bze-alphateam/cointrunk-discord-bridge/app/service"
	"github.com/spf13/cobra"
)

var roleChangerCmd = &cobra.Command{
	Use:   "role-changer",
	Short: "Starts Discord bot for role management",
	Long:  `Starts a Discord bot that allows users to subscribe/unsubscribe to get or remove a specific role`,
	Run: func(cmd *cobra.Command, args []string) {
		err := appCfg.Discord.ValidateRoleChanger()
		if err != nil {
			log.Fatal(err)
		}

		bot, err := service.NewDiscordBot(
			appCfg.Discord.RoleChangerToken,
			appCfg.Discord.AppID,
			appCfg.Discord.GuildID,
			appCfg.Discord.RoleID,
		)
		if err != nil {
			log.Fatal(err)
		}

		err = bot.Start()
		if err != nil {
			log.Fatal(err)
		}

		doneChan := make(chan bool, 1)
		listenSigTerm(func() {
			log.Println("Shutting down bot...")
			err := bot.Stop()
			if err != nil {
				log.Printf("Error stopping bot: %v", err)
			}
			doneChan <- true
		})

		<-doneChan
	},
}
