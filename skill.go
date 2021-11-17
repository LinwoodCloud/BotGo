package main

import "github.com/bwmarrin/discordgo"

type SkillUser struct {
	ID    string `gorm:"primarykey"`
	Skill string
}
type Skill struct {
	ID       string  `gorm:"primarykey"`
	Category *string `gorm:"default:''"`
	Name     string
	Link     string
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
)
