package finance

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/util"
)

type AccountValues struct {
	Name  string  `json:"name"`
	Value float64 `json:"total"`
}

func GetAccountValues() ([]AccountValues, error) {
	notion_db := config.Application.NotionAccountTableId
	dbItems, cliError := util.Notion.QueryDatabase(context.Background(), notion_db, nil)

	if cliError != nil {
		return []AccountValues{}, cliError
	}

	aggVal := make([]AccountValues, len(dbItems.Results))
	for i, r := range dbItems.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		aggVal[i] = AccountValues{
			props["Name"].Title[0].PlainText,
			*props["Total"].Rollup.Number,
		}
	}

	return aggVal, nil
}
