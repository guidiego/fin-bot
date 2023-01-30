package config

type Config struct {
	ApiPort                  string
	ApiToken                 string
	BankSlugToId             map[string]string
	BankIdToSlug             map[string]string
	Budgets                  map[string]string
	UsersToRemember          string `json:"UsersToRemember"`
	NotionScheduleTableId    string `json:"NotionScheduleTableId"`
	NotionTransactionTableId string `json:"NotionTransactionTableId"`
	NotionAccountTableId     string `json:"NotionAccountTableId"`
	NotionBudgetTableId      string `json:"NotionBudgetTableId"`
	DiscordBotToken          string
	DiscordTrackChannelId    string `json:"DiscordTrackChannelId"`
	DiscordRememberChannelId string `json:"DiscordRememberChannelId"`
	DiscordBalanceChannnelId string `json:"DiscordBalanceChannnelId"`
}

var (
	Application = Config{}
)
