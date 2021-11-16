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
		"hello": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hello, I'm a bot!",
				},
			})
		},
		"info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "**Information**",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{ID: "803156914028281856"},
									Label: "Programm Chest",
									Style: discordgo.LinkButton,
									URL:   "https://programm-chest.dev",
								},
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{ID: "834458109480140850"},
									Label: "Linwood",
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
