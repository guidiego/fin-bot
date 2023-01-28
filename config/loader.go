package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fin-bot/util"
)

func Load() error {
	notion_db := os.Getenv("NOTION_CONFIG_DB_ID")
	cRes, configErr := util.Notion.QueryDatabase(context.Background(), notion_db, nil)

	if configErr != nil {
		return configErr
	}

	cMap := map[string]string{}
	for _, r := range cRes.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		key := props["Key"].Title[0].PlainText
		value := props["Value"].RichText[0].PlainText
		cMap[key] = value
	}

	data, marshalErr := json.Marshal(cMap)
	if marshalErr != nil {
		return marshalErr
	}

	unmarshalErr := json.Unmarshal(data, &Application)
	if marshalErr != nil {
		return unmarshalErr
	}

	accountRes, accountErr := util.Notion.QueryDatabase(context.Background(), Application.NotionAccountTableId, nil)
	if accountErr != nil {
		return accountErr
	}

	bankSlugToId := map[string]string{}
	bankIdToSlug := map[string]string{}
	for _, r := range accountRes.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		name := props["Name"].Title[0].PlainText
		id := r.ID

		bankSlugToId[name] = id
		bankIdToSlug[id] = name
	}

	budgetRes, budgetErr := util.Notion.QueryDatabase(context.Background(), Application.NotionBudgetTableId, nil)
	if budgetErr != nil {
		return budgetErr
	}

	budgets := map[string]string{}
	for _, r := range budgetRes.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		name := props["Budget"].Title[0].PlainText
		id := r.ID

		budgets[id] = name
	}

	Application.BankIdToSlug = bankIdToSlug
	Application.BankSlugToId = bankSlugToId
	Application.Budgets = budgets

	Application.ApiPort = os.Getenv("PORT")
	Application.DiscordBotToken = os.Getenv("DISCORD_BOT_TOKEN")

	return nil
}
