package finance

import (
	"os"

	"github.com/dstotijn/go-notion"
)

var (
	notionCli = notion.NewClient(os.Getenv("NOTION_TOKEN"))
)
