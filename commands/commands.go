package commands

import "github.com/bwmarrin/discordgo"

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)
type SlashCommand struct {
	Spec    discordgo.ApplicationCommand
	Handler CommandHandler
}

var (
	Commands = map[string]*SlashCommand{
		insertTransactionCMD.Name: {
			insertTransactionCMD,
			insertTransactionHandler,
		},
	}
)
