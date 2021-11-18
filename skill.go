package main

import (
	"github.com/bwmarrin/discordgo"
)

type SkillUser struct {
	ID    string
	User  string `gorm:"primaryKey"`
	Skill string
}
type Skill struct {
	ID       string  `gorm:"primaryKey"`
	Category *string `gorm:"default:''"`
	Name     string
	Link     string
}
type SkillCategory struct {
	ID            string `gorm:"primaryKey"`
	Name          string
	Description   string
	Emoji         string
	IsCustomEmoji bool
}

func GetSkillsByUser(userID string) []*Skill {
	var skills []SkillUser
	database.Where("user = ?", userID).Find(&skills)
	var skillList []*Skill
	for _, skill := range skills {
		var skillData Skill
		database.Where("id = ?", skill.Skill).Find(&skillData)
		skillList = append(skillList, &skillData)
	}
	return skillList
}
func GetCategorySkills(category string) []*Skill {
	var skills []*Skill
	database.Where("category = ?", category).Find(&skills)
	return skills
}
func GetSkillCategoriesByUser(userID string) []string {
	var skills []SkillUser
	database.Where("user = ?", userID).Find(&skills)
	var skillCategories []string
	for _, skill := range skills {
		var skillData Skill
		database.Where("id = ?", skill.Skill).Find(&skillData)
		if *skillData.Category != "" {
			skillCategories = append(skillCategories, *skillData.Category)
		}
	}
	return skillCategories
}
func GetSkillsByUserAndCategory(userID string, category string) []*Skill {
	var skills []SkillUser
	database.Where("user = ? AND category = ?", userID, category).Find(&skills)
	var skillList []*Skill
	for _, skill := range skills {
		var skillData Skill
		database.Where("id = ?", skill.Skill).Find(&skillData)
		skillList = append(skillList, &skillData)
	}
	return skillList
}

func (su SkillUser) GetSkill() *Skill {
	s := Skill{}
	database.First(&s, su.Skill)
	return &s
}

func GetSkill(name string) *Skill {
	s := Skill{Name: name}
	database.First(&s, name)
	return &s
}

func SetupSkill() {
	database.AutoMigrate(&SkillUser{})
	database.AutoMigrate(&Skill{})
}

func EmbedSkills(skills []*Skill) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField
	skillByCategory := make(map[*string][]*Skill)
	for _, skill := range skills {
		if _, ok := skillByCategory[skill.Category]; !ok {
			skillByCategory[skill.Category] = make([]*Skill, 0)
		}
		skillByCategory[skill.Category] = append(skillByCategory[skill.Category], skill)
	}
	for category, skills := range skillByCategory {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**" + *category + "**",
			Value:  "",
			Inline: false,
		})
		for _, skill := range skills {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   skill.Name,
				Value:  skill.Link,
				Inline: true,
			})
		}
	}
	return fields
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
			Description: "Get all skills from an user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user to check the skills of.",
					Type:        discordgo.ApplicationCommandOptionUser,
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
				cateogries := GetSkillCategoriesByUser(user.ID)
				if len(cateogries) == 0 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "No skills found for user " + user.Username}})
					return
				}

				skills := GetSkillsByUserAndCategory(user.ID, cateogries[0])

				embed := &discordgo.MessageEmbed{
					Title:  "Skills",
					Color:  0x00ff00,
					Fields: EmbedSkills(skills),
				}
				options := make([]discordgo.SelectMenuOption, 0)
				for _, category := range cateogries {
					options = append(options, discordgo.SelectMenuOption{
						Label: category,
						Value: category,
					})
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: " ",
						Embeds:  []*discordgo.MessageEmbed{embed},
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										CustomID:    "skills",
										Placeholder: "Select a category",
										Options:     options,
									},
								},
							},
						},
					},
				})

			} else {
				cateogries := GetSkillCategoriesByUser(i.Member.User.ID)
				if len(cateogries) == 0 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "No skills found for user " + i.Member.User.Username}})
					return
				}

				skills := GetSkillsByUserAndCategory(i.Member.User.ID, cateogries[0])

				embed := &discordgo.MessageEmbed{
					Title:  "Skills",
					Color:  0x00ff00,
					Fields: EmbedSkills(skills),
				}
				options := make([]discordgo.SelectMenuOption, 0)
				for _, category := range cateogries {
					options = append(options, discordgo.SelectMenuOption{
						Label: category,
						Value: category,
					})
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: " ",
						Embeds:  []*discordgo.MessageEmbed{embed},
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										CustomID:    "skills",
										Placeholder: "Select a category",
										Options:     options,
									},
								},
							},
						},
					},
				})
			}
		},
	}
)
