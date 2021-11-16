package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	funCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "hello",
			Description: "Say hello",
		},
		{
			Name:        "info",
			Description: "Get information about Programm Chest",
		},
	}

	funCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hello, I'm a bot!",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{Name: "🌐"},
									Label: "Visit website",
									URL:   "https://programm-chest.dev",
									Style: discordgo.LinkButton,
								},
							},
						},
					},
				},
			})
		},
		"hello": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hello :D",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{ID: "834458109480140850"},
									Label: "Click me",
									Style: discordgo.LinkButton,
									URL:   "https://linwood.dev",
								},
							},
						},
					},
				},
			})
		},
	}
)
