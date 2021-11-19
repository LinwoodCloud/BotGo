package main

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
)

type EconomyUser struct {
	User       string `gorm:"primaryKey"`
	Team       string `gorm:"primaryKey"`
	Currency   EconomyCurrency
	CurrencyID int64 `gorm:"default:0;primaryKey"`
	Amount     int
}

type EconomyCurrency struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	PluralName  string
	Description string
	Tradeable   bool
	Emoji       sql.NullString
	CustomEmoji bool
}

func (e *EconomyUser) AddCoins(amount int) {
	e.Amount += amount
}
func (e *EconomyUser) RemoveCoins(amount int) {
	e.Amount -= amount
}
func (e *EconomyUser) Save() {
	database.Save(e)
}
func GetEconomyUser(userID string, currency int64) *EconomyUser {
	eu := EconomyUser{User: userID, Amount: 0}
	database.Where(&EconomyUser{User: userID, CurrencyID: currency}).First(&eu)
	return &eu
}
func GetEconomyUsers(userID string) *EconomyUser {
	eu := EconomyUser{User: userID, Amount: 0}
	database.Where(&EconomyUser{User: userID}).First(&eu)
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
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "remove",
					Description: "Remove coins from a user.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "The user that will have their coins removed",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "count",
							Description: "The count of points which will be removed",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set",
					Description: "Set the amount of coins a user has.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "The user that will have their coins set",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "count",
							Description: "The count of points which will be set",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "create",
					Description: "Create a new currency.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "name",
							Type:        discordgo.ApplicationCommandOptionString,
							Description: "The name of the currency.",
							Required:    true,
						},
						{
							Name:        "shortname",
							Type:        discordgo.ApplicationCommandOptionString,
							Description: "The short name of the currency.",
							Required:    false,
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
			/*if len(i.ApplicationCommandData().Options) > 0 && i.ApplicationCommandData().Options[0].UserValue(s) != nil {
							user := i.ApplicationCommandData().Options[0].UserValue(s)
							eus := GetEconomyUsers(user.ID)
							if len(eus) == 0 {
			                    s.ChannelMessageSend(i.ChannelID, "That user has no coins.")
			                    return
			                }

						} else {
							eu := GetEconomyUser(i.Member.User.ID)
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: fmt.Sprintf("You have %d coins.", eu.Amount),
								},
							})
						}*/
		},
		"coinsadmin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.ApplicationCommandData().Options[0].Name {
			case "add":
				/*if i.Member.Permissions&discordgo.PermissionManageServer == 0 {
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
				})*/
			}
		},
	}
)
