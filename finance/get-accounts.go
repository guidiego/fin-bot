package finance

import (
	"context"
	"fmt"

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
		prop := props["Total"]
		total := *prop.Rollup.Number
		field, err := util.Notion.FindPagePropertyByID(context.Background(), r.ID, prop.ID, nil)

		if err == nil {
			total = *field.PropertyItem.Rollup.Number
		} else {
			fmt.Printf("%e\n", err)
		}

		aggVal[i] = AccountValues{
			props["Name"].Title[0].PlainText,
			total,
		}
	}

	return aggVal, nil
}
