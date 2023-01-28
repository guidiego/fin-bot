package finance

import (
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
	"github.com/google/uuid"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/util"
)

func buildPagePayload(dbid string, value float64, uuid string, bankId string, desc string, budgetId string) (notion.CreatePageParams, error) {
	title := []notion.RichText{
		{
			Text: &notion.Text{Content: uuid},
		},
	}

	description := []notion.RichText{
		{
			Text: &notion.Text{Content: desc},
		},
	}

	account := []notion.Relation{
		{
			ID: bankId,
		},
	}

	budget := []notion.Relation{
		{
			ID: budgetId,
		},
	}

	emoji := util.GetTransactionEmoji(value)
	return notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   dbid,
		Icon: &notion.Icon{
			Type:  "emoji",
			Emoji: &emoji,
		},
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Ref": notion.DatabasePageProperty{
				Title: title,
			},
			"Desc": notion.DatabasePageProperty{
				RichText: description,
			},
			"Conta": notion.DatabasePageProperty{
				Relation: account,
			},
			"Valor": notion.DatabasePageProperty{
				Number: &value,
			},
			"Budget": {
				Relation: budget,
			},
		},
	}, nil
}

func NewTransaction(value float64, account string, content string, budgetId string) error {
	notion_db := config.Application.NotionTransactionTableId
	p, buildErr := buildPagePayload(notion_db, value, uuid.New().String(), account, content, budgetId)

	if buildErr != nil {
		return buildErr
	}

	_, cliError := util.Notion.CreatePage(context.Background(), p)

	if cliError != nil {
		fmt.Printf("%e", cliError)
		return cliError
	}

	return nil
}
