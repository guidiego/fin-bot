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
		if i.ApplicationCommandData().Name == "t" {
			commands.InsertTransactionHandler(s, i)
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
	cmd, err2 := s.ApplicationCommandCreate(s.State.User.ID, guildId, &commands.InsertTransactionCMD)

	if err2 != nil {
		log.Fatal(err2)
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	err3 := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
	if err3 != nil {
		log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
	}
}
