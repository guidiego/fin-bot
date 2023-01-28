package finance

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/util"
)

func GetScheduledTransactions() ([]notion.Page, error) {
	notion_db := config.Application.NotionScheduleTableId
	dbItems, cliError := util.Notion.QueryDatabase(context.Background(), notion_db, nil)

	if cliError != nil {
		return []notion.Page{}, cliError
	}

	return dbItems.Results, nil
}
