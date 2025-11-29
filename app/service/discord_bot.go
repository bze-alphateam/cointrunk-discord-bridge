package service

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	session *discordgo.Session
	guildID string
	roleID  string
	appID   string
}

func NewDiscordBot(token, appID, guildID, roleID string) (*DiscordBot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	return &DiscordBot{
		session: session,
		guildID: guildID,
		roleID:  roleID,
		appID:   appID,
	}, nil
}

func (b *DiscordBot) Start() error {
	b.session.AddHandler(b.handleInteraction)

	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord connection: %w", err)
	}

	log.Println("Discord bot is now running")

	err = b.registerCommands()
	if err != nil {
		return fmt.Errorf("error registering commands: %w", err)
	}

	return nil
}

func (b *DiscordBot) Stop() error {
	err := b.unregisterCommands()
	if err != nil {
		log.Printf("Error unregistering commands: %v", err)
	}

	return b.session.Close()
}

func (b *DiscordBot) registerCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "subscribe",
			Description: "Assigns you a role used for all the messages addressed to the entire community",
		},
		{
			Name:        "unsubscribe",
			Description: "Removes the role assigned to you for all the messages addressed to the entire community",
		},
	}

	for _, cmd := range commands {
		_, err := b.session.ApplicationCommandCreate(b.appID, b.guildID, cmd)
		if err != nil {
			return fmt.Errorf("cannot create '%s' command: %w", cmd.Name, err)
		}
		log.Printf("Registered command: %s", cmd.Name)
	}

	return nil
}

func (b *DiscordBot) unregisterCommands() error {
	commands, err := b.session.ApplicationCommands(b.appID, b.guildID)
	if err != nil {
		return fmt.Errorf("cannot fetch commands: %w", err)
	}

	for _, cmd := range commands {
		err := b.session.ApplicationCommandDelete(b.appID, b.guildID, cmd.ID)
		if err != nil {
			log.Printf("Cannot delete '%s' command: %v", cmd.Name, err)
		}
	}

	return nil
}

func (b *DiscordBot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	userID := i.Member.User.ID
	var response string
	var err error

	switch i.ApplicationCommandData().Name {
	case "subscribe":
		err = b.addRole(userID)
		if err != nil {
			response = fmt.Sprintf("Failed to add role: %s", err.Error())
			log.Printf("Error adding role to user %s: %v", userID, err)
		} else {
			response = "Successfully subscribed! Role has been added."
			log.Printf("Added role to user %s", userID)
		}
	case "unsubscribe":
		err = b.removeRole(userID)
		if err != nil {
			response = fmt.Sprintf("Failed to remove role: %s", err.Error())
			log.Printf("Error removing role from user %s: %v", userID, err)
		} else {
			response = "Successfully unsubscribed! Role has been removed."
			log.Printf("Removed role from user %s", userID)
		}
	default:
		response = "Unknown command"
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

func (b *DiscordBot) addRole(userID string) error {
	return b.session.GuildMemberRoleAdd(b.guildID, userID, b.roleID)
}

func (b *DiscordBot) removeRole(userID string) error {
	return b.session.GuildMemberRoleRemove(b.guildID, userID, b.roleID)
}
