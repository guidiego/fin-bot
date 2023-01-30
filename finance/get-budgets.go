package finance

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/util"
)

type Budget struct {
	Icon  string
	Name  string
	Limit float64
	Used  float64
	Free  float64
}

func GetBudgets() ([]Budget, error) {
	notion_db := config.Application.NotionBudgetTableId
	dbItems, cliError := util.Notion.QueryDatabase(context.Background(), notion_db, nil)

	if cliError != nil {
		return []Budget{}, cliError
	}

	aggVal := make([]Budget, len(dbItems.Results))
	for i, r := range dbItems.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		aggVal[i] = Budget{
			*r.Icon.Emoji,
			props["Budget"].Title[0].PlainText,
			*props["Limit"].Number,
			*props["Used"].Rollup.Number,
			*props["Free"].Formula.Number,
		}
	}

	return aggVal, nil
}
