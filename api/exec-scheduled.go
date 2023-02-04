package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dstotijn/go-notion"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/finance"
	"github.com/guidiego/fin-bot/util"
)

func ExecScheduleRoute(s *discordgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsTokenOk(w, r) {
			return
		}

		results, err := finance.GetScheduledTransactions()
		usersToAlert := config.Application.UsersToRemember
		trackChannelId := config.Application.DiscordTrackChannelId
		rememberChannelId := config.Application.DiscordRememberChannelId
		now := time.Now()
		today := float64(now.Day())

		if err != nil {
			log.Printf("%e", err)
		}

		for _, r := range results {
			props := r.Properties.(notion.DatabasePageProperties)
			day := *props["Day"].Number

			if today != day {
				continue
			}

			value := *props["Value"].Number
			ref := props["Ref"].Title[0].Text.Content
			automatic := *props["AutoDebit"].Checkbox

			rawMsg := ""
			msg := ""
			channelId := rememberChannelId

			if !automatic {
				rawMsg = "%s **[%s] %.2f€** precisa ser feito hoje! CC: %s"
				msg = fmt.Sprintf(rawMsg, util.GetTransactionEmoji(value), ref, value, usersToAlert)
			} else {
				bankId := props["Account"].Relation[0].ID
				budgetRel := props["Budget"].Relation
				budget := ""

				if len(budgetRel) > 0 {
					budget = budgetRel[0].ID
				}

				err := finance.NewTransaction(value, bankId, ref, budget)

				if err != nil {
					fmt.Printf("%e", err)
					continue
				}

				channelId = trackChannelId
				rawMsg = "%s **%.2f€** transacionado em **%s** - %s, CC: %s"
				msg = fmt.Sprintf(rawMsg, util.GetTransactionEmoji(value), value, config.Application.BankIdToSlug[bankId], ref, usersToAlert)
			}

			s.ChannelMessageSend(channelId, msg)
		}

		fmt.Fprint(w, "Done")
	}
}
