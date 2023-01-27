package util

func GetTransactionEmoji(value float64) string {
	emoji := "🟢"

	if value < 0 {
		emoji = "🔴"
	}

	return emoji
}
