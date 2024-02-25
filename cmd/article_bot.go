package cmd

import (
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/command"
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/factory"
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/repository"
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/service"
	"github.com/spf13/cobra"
	"log"
)

var articleBotCmd = &cobra.Command{
	Use:   "start",
	Short: "Publishes latest articles on discord",
	Long:  `Publishes latest CoinTrunk.io articles from BZE Blockchain to your discord channel`,
	Run: func(cmd *cobra.Command, args []string) {
		err := appCfg.Validate()
		if err != nil {
			log.Fatal(err)
		}

		repo, err := repository.NewHistory(appCfg.History.FileName, appCfg.History.FilePath)
		if err != nil {
			log.Fatal(err)
		}

		dp, err := service.NewRestDataProvider(appCfg.RestUrl)
		if err != nil {
			log.Fatal(err)
		}

		msgFactory, err := factory.NewWebhookFactory(dp)
		if err != nil {
			log.Fatal(err)
		}

		tm, err := service.NewTendermintService(appCfg.RpcUrl)
		if err != nil {
			log.Fatal(err)
		}

		webhook, err := service.NewDiscordWebhook(appCfg.Discord.Webhook)
		if err != nil {
			log.Fatal(err)
		}

		appCmd, err := command.NewDiscordArticle(dp, repo, msgFactory, webhook, tm)
		if err != nil {
			log.Fatal(err)
		}

		doneChan := make(chan bool, 1)
		listenSigTerm(func() {
			doneChan <- true
		})

		go func() {
			err = appCmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			doneChan <- true
		}()

		<-doneChan
	},
}
