package finance

import (
	"context"
	"errors"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/google/uuid"
	"github.com/guidiego/fin-bot/util"
)

var bankIds = map[string]string{
	"commerz":  "688b791acb334cea8526f25ae9c116f6",
	"amex":     "e460bf056c8546f9add452e67d2b833c",
	"dinheiro": "bbc28caaefc8470f9d824e8139febd3d",
	"crypto":   "181a341facaf4939bae5392de14d5f36",
}

func buildPagePayload(dbid string, value float64, uuid string, bankSlug string, desc string) (notion.CreatePageParams, error) {
	bankId, bankIdExists := bankIds[bankSlug]

	if !bankIdExists {
		return notion.CreatePageParams{}, errors.New("bank slug not supported")
	}

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

	relation := []notion.Relation{
		{
			ID: bankId,
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
				Relation: relation,
			},
			"Valor": notion.DatabasePageProperty{
				Number: &value,
			},
		},
	}, nil
}

func NewTransaction(value float64, account string, content string) error {
	notion_db := os.Getenv("NOTION_TRANSACTION_DB_ID")
	p, buildErr := buildPagePayload(notion_db, value, uuid.New().String(), account, content)

	if buildErr != nil {
		return buildErr
	}

	_, cliError := notionCli.CreatePage(context.Background(), p)

	if cliError != nil {
		return cliError
	}

	return nil
}
