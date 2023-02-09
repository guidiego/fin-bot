package finance

import (
	"context"
	"fmt"
	"time"

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

	emoji := util.GetTransactionEmoji(value)
	database_properties := notion.DatabasePageProperties{
		"Ref": notion.DatabasePageProperty{
			Title: title,
		},
		"Desc": notion.DatabasePageProperty{
			RichText: description,
		},
		"Account": notion.DatabasePageProperty{
			Relation: account,
		},
		"Value": notion.DatabasePageProperty{
			Number: &value,
		},
	}

	if budgetId != "" {
		database_properties["Budget"] = notion.DatabasePageProperty{
			Relation: []notion.Relation{
				{
					ID: budgetId,
				},
			},
		}
	}
	return notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   dbid,
		Icon: &notion.Icon{
			Type:  "emoji",
			Emoji: &emoji,
		},
		DatabasePageProperties: &database_properties,
	}, nil
}

func NewTransaction(value float64, account string, content string, budgetId string) error {
	notion_db := config.Application.NotionTransactionTableId
	t := time.Now()
	token := t.Format("2006-01")
	id := token + "-" + uuid.New().String()
	p, buildErr := buildPagePayload(notion_db, value, id, account, content, budgetId)

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
