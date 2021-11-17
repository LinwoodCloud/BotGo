package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type EconomyUser struct {
	ID    string `gorm:"primarykey"`
	Coins int
}

func (e *EconomyUser) AddCoins(amount int) {
	e.Coins += amount
}
func (e *EconomyUser) RemoveCoins(amount int) {
	e.Coins -= amount
}
func (e *EconomyUser) Save() {
	database.Save(e)
}
func GetEconomyUser(userID string) *EconomyUser {
	eu := EconomyUser{ID: userID, Coins: 0}
	database.FirstOrCreate(&eu, userID)
	return &eu
}

func setupEconomy() {
	database.AutoMigrate(&EconomyUser{})
}

var (
	economyCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "coins",
			Description: "Displays the amount of coins the user has.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to check the coins of.",
					Required:    false,
				},
			},
		},
		{
			Name:        "addcoins",
			Description: "Add coins to the given user",

			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user that will be given the coins",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "count",
					Description: "The count of points which will be given",
					Required:    true,
				},
			},
		},
		{
			Name:        "take",
			Description: "Takes the specified user the specified amount of coins.",
		},
	}
	economyCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"coins": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(i.ApplicationCommandData().Options) > 0 && i.ApplicationCommandData().Options[0].UserValue(s) != nil {
				user := i.ApplicationCommandData().Options[0].UserValue(s)
				eu := GetEconomyUser(user.ID)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("%s has %d coins.", user.Username, eu.Coins),
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
		"addcoins": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				// Here we need to convert raw interface{} value to wanted type.
				// Also, as you can see, here is used utility functions to convert the value
				// to particular type. Yeah, you can use just switch type,
				// but this is much simpler
				i.ApplicationCommandData().Options[0].UserValue(s).ID,
				i.ApplicationCommandData().Options[1].IntValue(),
			}
			if i.Member.Permissions&discordgo.PermissionManageServer == 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   1 << 6,
						Content: "You don't have permission to do that.",
					},
				})
			}
			if len(margs) != 2 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   1 << 6,
						Content: "Invalid arguments.",
					},
				})
				return
			}
			// Get name of user
			user, err := s.User(margs[0].(string))
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "User not found.",
					},
				})
				return
			}
			eu := GetEconomyUser(margs[0].(string))

			eu.AddCoins(int(margs[1].(int64)))
			eu.Save()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You have added %d coins to %s.", margs[1].(int64), user.Username),
				},
			})
		},
	}
)
