package commands

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/guidiego/fin-bot/finance"
	"github.com/guidiego/fin-bot/util"
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

func InteractionRespond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func insertTransactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = *opt
	}

	valueOpt, isValueOk := optionMap["value"]
	accountOpt, isAccountOk := optionMap["account"]

	if !isValueOk || !isAccountOk {
		InteractionRespond(s, i, "[Err] Missing Parameters")
		return
	}

	content := ""
	contentOpt, isContentOk := optionMap["description"]

	if isContentOk {
		content = contentOpt.StringValue()
	}

	transactionVal, fParseErr := strconv.ParseFloat(valueOpt.StringValue(), 64)
	if fParseErr != nil {
		InteractionRespond(s, i, "[Err] Invalid Value")
		return
	}

	account := accountOpt.StringValue()
	err := finance.NewTransaction(
		transactionVal, account, content,
	)

	if err != nil {
		fmt.Printf("err: %e\n", err)
		InteractionRespond(s, i, "[Err] Problems with notion API!")
		return
	}

	emoji := util.GetTransactionEmoji(transactionVal)
	InteractionRespond(s, i,
		fmt.Sprintf("%s **%.2f€** transacionado em **%s**", emoji, transactionVal, account),
	)
}
