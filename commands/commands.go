package commands

import "github.com/bwmarrin/discordgo"

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)
type SlashCommand struct {
	Spec    discordgo.ApplicationCommand
	Handler CommandHandler
}

func Build() map[string]*SlashCommand {
	return map[string]*SlashCommand{
		insertTransactionCmdName: {
			insertTransactionCommandBuild(),
			insertTransactionHandler,
		},
	}
}
