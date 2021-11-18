package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

var (
	adminCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "userinfo",
			Description: "Get information about a user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to get information about.",
				},
			},
		},
	}
	adminCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"userinfo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			member := i.Member
			user := member.User
			if len(i.ApplicationCommandData().Options) == 1 {
				user = i.ApplicationCommandData().Options[0].UserValue(s)
			}
			joined, err := member.JoinedAt.Parse()
			if err != nil {
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: " ",
					Embeds: []*discordgo.MessageEmbed{{
						Title: "User Information",
						Image: &discordgo.MessageEmbedImage{
							URL: user.AvatarURL(""),
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Username",
								Value:  user.Username,
								Inline: true,
							},
							{
								Name:   "ID",
								Value:  user.ID,
								Inline: true,
							},
							{
								Name:   "Discriminator",
								Value:  user.Discriminator,
								Inline: true,
							},
							{
								Name:   "Joined",
								Value:  joined.Format(time.RFC3339),
								Inline: true,
							},
						},
					},
					},
				},
			})
		},
	}
)
