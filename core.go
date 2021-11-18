package main

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
)

type CoreUser struct {
	ID     string `gorm:"primarykey"`
	Locale sql.NullString
}

func (cu *CoreUser) GetLocaleOrDefault() string {
	if cu.Locale.Valid {
		return cu.Locale.String
	}
	user, err := s.User(cu.ID)
	if err != nil {
		return "en-US"
	}
	return user.Locale
}

func (cu *CoreUser) Save() {
	database.Save(cu)
}

// GetUser returns the user from the database
func GetUser(userID string) *CoreUser {
	cu := CoreUser{ID: userID}
	database.FirstOrCreate(&cu, userID)
	return &cu
}
func SetupCore() {
	database.AutoMigrate(&CoreUser{})
}

var (
	coreCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "locale",
			Description: "Get your current locale",
		},
		{
			Name:        "setlocale",
			Description: "Set your locale",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "locale",
					Description: "Your preferred locale. Set nothing to reset it",
					Required:    false,
				},
			},
		},
	}
	coreCommandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"locale": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			cu := GetUser(i.User.ID)
			if cu.Locale.Valid {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Your locale is `" + cu.Locale.String + "`",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You don't have a custom locale set. Using the default `" + i.Member.User.Locale + "`",
					},
				})
			}
		},
		"setlocale": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var locale string
			if len(i.ApplicationCommandData().Options) == 1 || i.ApplicationCommandData().Options[0].Value != "" {
				locale = i.ApplicationCommandData().Options[0].Value.(string)
			}
			cu := GetUser(i.User.ID)
			cu.Locale.String = locale
			cu.Save()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your locale has been set to " + locale,
				},
			})
		},
	}
)
