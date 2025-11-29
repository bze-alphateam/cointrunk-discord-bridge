package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bze-alphateam/cointrunk-discord-bridge/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appCfg = &config.App{}

var rootCmd = &cobra.Command{
	Use:   "ctbot",
	Short: "CoinTrunk Discord Bot is publishing CoinTrunk.io articles to your discord",
	Long: `A discord bot listening to BZE blockchain for new CoinTrunk.io articles and publish them to a discord 
			channel of your choice`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func init() {
	rootCmd.AddCommand(articleBotCmd)
	rootCmd.AddCommand(roleChangerCmd)
}

func initConfig() {
	viper.SetConfigName("ctbot") // name of config file (without extension)
	viper.SetConfigType("yaml")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")     // optionally look for config in the working directory
	err := viper.ReadInConfig()  // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	err = viper.Unmarshal(appCfg)
	//force setting these two keys as viper fails to unmarshal them correctly.
	//odd bug ? or something... to avoid wasting time studying unmarshal and viper use this workaround
	appCfg.RestUrl = viper.GetString("rest_url")
	appCfg.RpcUrl = viper.GetString("rpc_url")
	appCfg.Discord.RoleChangerToken = viper.GetString("discord.role_changer_token")
	appCfg.Discord.AppID = viper.GetString("discord.app_id")
	appCfg.Discord.GuildID = viper.GetString("discord.guild_id")
	appCfg.Discord.RoleID = viper.GetString("discord.role_id")
	if err != nil {
		log.Fatalf("Unable to unmarshal config into struct, %v", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func listenSigTerm(cancel func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		// Listen for the interrupt signal.
		<-sigChan
		fmt.Println("\nShutdown signal received, exiting...")
		cancel() // Cancel the context.
	}()
}
