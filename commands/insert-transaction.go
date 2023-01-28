package commands

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/finance"
	"github.com/guidiego/fin-bot/util"
)

func InteractionRespond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

var insertTransactionCmdName = "t"

func insertTransactionCommandBuild() discordgo.ApplicationCommand {
	accountChoices := []*discordgo.ApplicationCommandOptionChoice{}
	budgetChoices := []*discordgo.ApplicationCommandOptionChoice{}

	for id, name := range config.Application.BankIdToSlug {
		accountChoices = append(accountChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: id,
		})
	}

	for id, name := range config.Application.Budgets {
		budgetChoices = append(budgetChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: id,
		})
	}

	return discordgo.ApplicationCommand{
		Name:        insertTransactionCmdName,
		Description: "This command inserts a transaction in your Notion Page",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "account",
				Description: "Source account for transaction",
				Required:    true,
				Choices:     accountChoices,
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
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "budget",
				Description: "Choose the budget from this value",
				Choices:     budgetChoices,
			},
		},
	}
}
func insertTransactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = *opt
	}

	valueOpt, isValueOk := optionMap["value"]
	actionnOpt, isActionnOk := optionMap["action"]
	accountOpt, isAccountOk := optionMap["account"]

	if !isValueOk || !isAccountOk || !isActionnOk {
		InteractionRespond(s, i, "[Err] Missing Parameters")
		return
	}

	desc := ""
	content := ""
	budget := ""
	contentOpt, isContentOk := optionMap["description"]
	budgetOpt, isBudgetOk := optionMap["budget"]

	if isContentOk {
		content = contentOpt.StringValue()
		desc = " - " + content
	}

	if isBudgetOk {
		budget = budgetOpt.StringValue()
	}

	transactionVal, fParseErr := strconv.ParseFloat(valueOpt.StringValue(), 64)

	if (actionnOpt.StringValue() == "rm" && transactionVal > 0) || (actionnOpt.StringValue() == "add" && transactionVal < 0) {
		transactionVal = transactionVal * -1
	}

	if fParseErr != nil {
		InteractionRespond(s, i, "[Err] Invalid Value")
		return
	}

	account := accountOpt.StringValue()
	fmt.Println(account)
	err := finance.NewTransaction(
		transactionVal, account, content, budget,
	)

	if err != nil {
		fmt.Printf("%e \n", err)
		InteractionRespond(s, i, "[Err] Problems with notion API!")
		return
	}

	emoji := util.GetTransactionEmoji(transactionVal)
	InteractionRespond(s, i,
		fmt.Sprintf("%s **%.2f€** transacionado em **%s**%s", emoji, transactionVal, config.Application.BankIdToSlug[account], desc),
	)
}
