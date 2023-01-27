package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	insertTransactionCMD = discordgo.ApplicationCommand{
		Name:        "t",
		Description: "This command inserts a transaction in your Notion Page",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "account",
				Description: "Source account for transaction",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "commerz",
						Value: "commerz",
					},
					{
						Name:  "amex",
						Value: "amex",
					},
					{
						Name:  "dinheiro",
						Value: "dinheiro",
					},
					{
						Name:  "crypto",
						Value: "crypto",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "action",
				Description: "Remove or add money",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "add",
						Value: "add",
					},
					{
						Name:  "rm",
						Value: "rm",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "value",
				Description: "value of transaction",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "description",
				Description: "description of the transaction",
			},
		},
	}
)

func insertTransactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = *opt
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"bank %s",
				optionMap["account"].Value,
			),
		},
	})
}
