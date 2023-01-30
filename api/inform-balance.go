package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/finance"
)

func getAccountMessage() (string, error) {
	msg := "**🏦 Accounts**\n\n"
	accounts, err := finance.GetAccountValues()

	if err != nil {
		return "", err
	}

	total := 0.0
	for _, v := range accounts {
		total = total + v.Value
		msg = msg + fmt.Sprintf(">  **%s:**    %.2f€\n", v.Name, v.Value)
	}

	msg = msg + fmt.Sprintf(">  **TOTAL:    %.2f**€\n\n\n\n", total)
	return msg, nil
}

func getBudgetsMessage() (string, error) {
	msg := "**👛 Budgets**\n\n"
	budgets, err := finance.GetBudgets()

	if err != nil {
		return "", err
	}

	total := len(budgets)
	for i, v := range budgets {
		used := -1 * v.Used
		emoji := "🟢"
		template := ">    %s **(%.2f€/%.2f€)**   restantes em   **%s %s**   _(livre: %.2f€)_\n"

		if i < total-1 {
			template = template + "> \n"
		}

		if v.Free < (v.Limit / 3) {
			emoji = "🔴"
		} else if v.Free < (v.Limit / 2) {
			emoji = "🟡"
		}

		template += " "
		msg = msg + fmt.Sprintf(template, emoji, used, v.Limit, v.Icon, v.Name, v.Free)
	}

	return msg, nil
}

func InformBalanceRoute(s *discordgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsTokenOk(w, r) {
			return
		}

		date := time.Now()
		msg := "🤑  **IT'S TIME TO YOUR BALANCE!** 🤑\n"
		msg = msg + fmt.Sprintf("📅  _%s_\n\n\n", date.Format("02/01/2006"))
		aMsg, aErr := getAccountMessage()

		if aErr != nil {
			fmt.Fprintf(w, "%e", aErr)
			return
		}

		bMsg, bErr := getBudgetsMessage()

		if bErr != nil {
			fmt.Fprintf(w, "%e", bErr)
			return
		}

		msg = msg + aMsg + bMsg + fmt.Sprintf("\n\n CC: %s\n", config.Application.UsersToRemember)

		s.ChannelMessageSend(config.Application.DiscordBalanceChannnelId, msg)

		fmt.Fprint(w, "Done")
	}
}
