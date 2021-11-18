package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type EconomyUser struct {
	ID    string `gorm:"primarykey"`
	Coins int
}

type EconomyCurrency struct {
	ID          string `gorm:"primarykey"`
	Name        string
	ShortName   string
	Description string
	Tradeable   bool
	Icon        string
	CustomIcon  bool
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

func SetupEconomy() {
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
			Name:        "coinsadmin",
			Description: "Manage coins. Only usable by the admins.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "add",
					Description: "Add coins to a user.",
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
			},
		},
		{
			Name:        "givecoins",
			Description: "Transfer your coins to another user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user that will receive the coins",
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
		"coinsadmin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.ApplicationCommandData().Options[0].Name {
			case "add":
				if i.Member.Permissions&discordgo.PermissionManageServer == 0 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Flags:   1 << 6,
							Content: "You don't have permission to do that.",
						},
					})
				}
				// Get name of user
				user := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
				eu := GetEconomyUser(user.ID)

				count := int(i.ApplicationCommandData().Options[0].Options[1].IntValue())
				eu.AddCoins(count)
				eu.Save()
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("You have added %d coins to %s.", count, user.Username),
					},
				})
			}
		},
	}
)
