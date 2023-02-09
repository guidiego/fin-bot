package finance

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/util"
)

type Goal struct {
	Name      string
	Desired   float64
	AlreadyIn float64
}

func GetGoals() ([]Goal, error) {
	notion_db := config.Application.NotionGoalTableId
	dbItems, cliError := util.Notion.QueryDatabase(context.Background(), notion_db, nil)

	if cliError != nil {
		return []Goal{}, cliError
	}

	aggVal := make([]Goal, len(dbItems.Results))
	for i, r := range dbItems.Results {
		props := r.Properties.(notion.DatabasePageProperties)
		aggVal[i] = Goal{
			props["Name"].Title[0].PlainText,
			*props["Desired"].Number,
			*props["AlreadyIn"].Number,
		}
	}

	return aggVal, nil
}
