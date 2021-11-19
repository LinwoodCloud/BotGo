package main

import (
	"errors"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
)

var s *discordgo.Session
var database *gorm.DB

func init() { flag.Parse() }

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func init() {
	if exists(".env") {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	var err error
	token := os.Getenv("CHEST_BOT_TOKEN")
	if token == "" {
		log.Fatal("CHEST_BOT_TOKEN not set")
	}
	s, err = discordgo.New("Bot " + token)
	database = buildDatabase()
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	GuildID = ""
	Cleanup = flag.Bool("cleanup", true, "Cleanup of commands")
)

func commands() []*discordgo.ApplicationCommand {
	var commands []*discordgo.ApplicationCommand
	commands = append(commands, adminCommands...)
	commands = append(commands, funCommands...)
	commands = append(commands, economyCommands...)
	commands = append(commands, skillCommands...)
	commands = append(commands, teamCommands...)
	return commands
}

func init() {
	ms := []map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		adminCommandHandlers,
		economyCommandHandlers,
		funCommandHandlers,
		skillCommandHandlers,
		teamCommandHandlers,
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
}
func Setup() {
	SetupCore()
	SetupEconomy()
	SetupSkill()
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	Setup()

	commands := commands()
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")

	if !*Cleanup {
		return
	}

	cmds, err := s.ApplicationCommands(s.State.User.ID, GuildID)

	for cmdId := range cmds {
		cmd := cmds[cmdId]
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, cmd.ID)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", cmd.Name, err)
		}
	}
}
