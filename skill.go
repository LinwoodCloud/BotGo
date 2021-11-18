package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type SkillUser struct {
	ID    string
	User  string `gorm:"primarykey"`
	Skill string
}
type Skill struct {
	ID       string  `gorm:"primarykey"`
	Category *string `gorm:"default:''"`
	Name     string
	Link     string
}

func GetSkillUser(userID string) *SkillUser {
	su := SkillUser{ID: userID}
	database.FirstOrCreate(&su, userID)
	return &su
}

func GetSkill(name string) *Skill {
	s := Skill{Name: name}
	database.First(&s, name)
	return &s
}

func setupSkill() {
	database.AutoMigrate(&SkillUser{})
}

var (
	skillCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "skill",
			Description: "Get skill list or the help to a specific ",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "skill",
					Description: "Skill name",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
				{
					Name:        "category",
					Description: "Skill category",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
		{
			Name:        "skills",
			Description: "Get all skills",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user to check the skills of.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
		{
			Name:        "addskill",
			Description: "Add a skill",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "skill",
					Description: "The skill to add.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "removeskill",
			Description: "Remove a skill",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "skill",
					Description: "The skill to remove.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "createskill",
			Description: "Create a skill",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "skill",
					Description: "The skill to create.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "category",
					Description: "The category of the skill.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
				{
					Name:        "link",
					Description: "The link to the skill.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
		{
			Name:        "editskill",
			Description: "Edit a skill",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "skill",
					Description: "The skill to edit.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "category",
					Description: "The category of the skill.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
				{
					Name:        "link",
					Description: "The link to the skill.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
		{
			Name:        "deleteskill",
			Description: "Delete a skill",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "skill",
					Description: "The skill to delete.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}
	skillCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"skills": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(i.ApplicationCommandData().Options) > 0 && i.ApplicationCommandData().Options[0].UserValue(s) != nil {
				user := i.ApplicationCommandData().Options[0].UserValue(s)
				su := GetSkillUser(user.ID)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("%s has %d coins.", user.Username, su.Skill),
					},
				})
			} else {
				eu := GetEconomyUser(i.Member.User.ID)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("You have %d coins.", eu.Coins),
					},
				})
			}
		},
	}
)
