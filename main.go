package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/guidiego/fin-bot/api"
	"github.com/guidiego/fin-bot/commands"
	"github.com/guidiego/fin-bot/config"
	"github.com/guidiego/fin-bot/util"
)

var s *discordgo.Session
var cmds map[string]*commands.SlashCommand

func init() {
	var err error

	util.InitNotion()
	err = config.Load()

	if err != nil {
		log.Fatalf("Fail Config: %v", err)
	}

	s, err = discordgo.New("Bot " + config.Application.DiscordBotToken)

	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	// cmds = commands.Build()
	// s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	cmd := cmds[i.ApplicationCommandData().Name]
	// 	if cmd != nil {
	// 		cmd.Handler(s, i)
	// 	}
	// })
}

func main() {
	// s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
	// 	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	// })
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// log.Println("Adding commands...")

	// registeredCommands := map[string]discordgo.ApplicationCommand{}
	// for _, v := range cmds {
	// 	cmd, creationErr := s.ApplicationCommandCreate(s.State.User.ID, "", &v.Spec)
	// 	registeredCommands[v.Spec.Name] = *cmd

	// 	if creationErr != nil {
	// 		log.Fatal(creationErr)
	// 	}
	// }

	http.HandleFunc("/exec-schedule", api.ExecScheduleRoute(s))
	http.HandleFunc("/inform-balance", api.InformBalanceRoute(s))

	if config.Application.ApiPort == "" {
		config.Application.ApiPort = "8080"
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	http.ListenAndServe(":"+config.Application.ApiPort, nil)
	<-stop

	// for _, v := range registeredCommands {
	// 	err3 := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
	// 	if err3 != nil {
	// 		log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
	// 	}
	// }
}
