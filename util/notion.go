package util

import (
	"os"

	"github.com/dstotijn/go-notion"
)

var Notion *notion.Client

func InitNotion() {
	Notion = notion.NewClient(os.Getenv("NOTION_TOKEN"))
}
