package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
)

var s *discordgo.Session
var database *gorm.DB

func init() { flag.Parse() }

var (
	GuildID = ""
)

func commands() []*discordgo.ApplicationCommand {
	var commands []*discordgo.ApplicationCommand
	commands = append(commands, adminCommands...)
	commands = append(commands, funCommands...)
	commands = append(commands, economyCommands...)
	return commands
}

func setup() {
	setupCore()
	setupEconomy()
}

func main() {
	var err error
	s, err = discordgo.New("Bot " + os.Getenv("CHEST_BOT_TOKEN"))
	database = buildDatabase()
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	setup()

	commands := commands()
	cmdIDs := make(map[string]string, len(commands))
	for _, v := range commands {
		rcmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		cmdIDs[rcmd.ID] = rcmd.Name
	}

	ms := []map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		adminCommandHandlers,
		economyCommandHandlers,
		funCommandHandlers,
	}
	res := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := res[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")

	for id, name := range cmdIDs {
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
}
