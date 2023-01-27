package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/guidiego/fin-bot/commands"
)

var s *discordgo.Session

func init() {
	var err error
	token := os.Getenv("DISCORD_BOT_TOKEN")
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmd := commands.Commands[i.ApplicationCommandData().Name]
		if cmd != nil {
			cmd.Handler(s, i)
		}
	})
}

func main() {
	guildId := os.Getenv("DISCORD_CHANNEL_ID")
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := map[string]discordgo.ApplicationCommand{}
	for _, v := range commands.Commands {
		cmd, creationErr := s.ApplicationCommandCreate(s.State.User.ID, guildId, &v.Spec)
		registeredCommands[v.Spec.Name] = *cmd

		if creationErr != nil {
			log.Fatal(creationErr)
		}
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	for _, v := range registeredCommands {
		err3 := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err3 != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
