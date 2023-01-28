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

var bankSlugs = map[string]string{
	"688b791acb334cea8526f25ae9c116f6": "commerz",
	"e460bf056c8546f9add452e67d2b833c": "amex",
	"bbc28caaefc8470f9d824e8139febd3d": "dinheiro",
	"181a341facaf4939bae5392de14d5f36": "crypto",
}

func ExecScheduleRoute(s *discordgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := finance.GetScheduledTransactions()
		today := time.Now()
		usersToAlert := config.Application.UsersToRemember
		// trackChannelId := config.Application.DiscordTrackChannelId
		rememberChannelId := config.Application.DiscordRememberChannelId

		if err != nil {
			log.Printf("%e", err)
		}

		for _, r := range results {
			props := r.Properties.(notion.DatabasePageProperties)

			dia := *props["Dia"].Number

			if float64(today.Day()) != dia {
				continue
			}

			value := *props["Valor"].Number
			ref := props["Ref"].Title[0].Text.Content
			automatic := *props["Debitar Automatico"].Checkbox

			rawMsg := ""
			msg := ""
			channelId := rememberChannelId
			if !automatic {
				rawMsg = "%s **[%s] %.2f€** precisa ser feito hoje! CC: %s"
				msg = fmt.Sprintf(rawMsg, util.GetTransactionEmoji(value), ref, value, usersToAlert)
			} else {
				bankId := props["Conta"].Relation[0].ID
				// err := finance.NewTransaction(value, bankSlugs[strings.Replace(bankId, "-", "", 5)], ref)

				// if err != nil {
				// 	fmt.Printf("%e", err)
				// 	continue
				// }

				// channelId = trackChannelId
				rawMsg = "%s **%.2f€** transacionado em **%s**, CC: %s"
				msg = fmt.Sprintf(rawMsg, util.GetTransactionEmoji(value), value, bankSlugs[bankId], usersToAlert)
			}

			s.ChannelMessageSend(channelId, msg)
		}

		fmt.Fprint(w, "Done")
	}
}
