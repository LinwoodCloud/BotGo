package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("CHEST_BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
}
