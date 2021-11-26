package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"strings"
)

type Team struct {
	Name        string `gorm:"primary_key"`
	Description string
	Members     []TeamMember `gorm:"foreignkey:TeamName;association_foreignkey:Name"`
}

type TeamMember struct {
	TeamName string `gorm:"primary_key"`
	Guild    string `gorm:"primary_key"`
	Role     TeamMemberRole
}
type TeamMemberRole int

const (
	TeamMemberRoleOwner TeamMemberRole = iota
	TeamMemberRoleModerator
	TeamMemberRoleMember
)

func (t *Team) AddMember(guild string) {
	t.Members = append(t.Members, TeamMember{Guild: guild, Role: TeamMemberRoleMember})
}

func (t *Team) RemoveMember(member string) {
	for i, m := range t.Members {
		if m.Guild == member {
			t.Members = append(t.Members[:i], t.Members[i+1:]...)
		}
	}
}

func (t *Team) GetMemberNames() []string {
	var tms []TeamMember
	query := database.Where("team_name = ?", t.Name).Find(&tms)
	if query.Error != nil {
		return []string{}
	}
	names := make([]string, len(tms))
	for i, m := range tms {
		guild, err := s.Guild(m.Guild)
		if err != nil {
			continue
		}
		names[i] = m.Role.GetEmoji() + " " + guild.Name + " (" + guild.ID + ")"
	}
	return names
}

func (t *Team) GetMember(member string) *TeamMember {
	for _, m := range t.Members {
		if m.Guild == member {
			return &m
		}
	}
	return &TeamMember{}
}

func (t *Team) BuildEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       t.Name,
		Description: t.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Members",
				// Value are team members separated with new line
				Value: strings.Join(t.GetMemberNames(), "\n"),
			},
		},
	}
}

func (t *TeamMember) Promote() {
	t.Role = TeamMemberRoleModerator
}
func (t *TeamMember) Demote() {
	if t.Role == TeamMemberRoleModerator {
		t.Role = TeamMemberRoleMember
	}
}
func (t *TeamMember) GetTeam() *Team {
	var team Team
	database.First(&team, "name = ?", t.TeamName)
	return &team
}

func (t *TeamMemberRole) String() string {
	switch *t {
	case TeamMemberRoleOwner:
		return "Owner"
	case TeamMemberRoleModerator:
		return "Moderator"
	case TeamMemberRoleMember:
		return "Member"
	default:
		return "Unknown"
	}
}
func (t *TeamMemberRole) GetEmoji() string {
	switch *t {
	case TeamMemberRoleOwner:
		return ":crown:"
	case TeamMemberRoleModerator:
		return ":star:"
	case TeamMemberRoleMember:
		return ":busts_in_silhouette:"
	default:
		return ":question:"
	}
}

func (t *Team) Save() {
	database.Save(t)
}

// CreateTeam creates a new team. Returns nil if the team already exists.
func CreateTeam(guildID string, name string, description string) *Team {
	team := Team{Name: name, Description: description, Members: []TeamMember{
		{Guild: guildID, Role: TeamMemberRoleOwner},
	}}
	if errors.Is(database.First(&team).Error, gorm.ErrRecordNotFound) {
		database.Create(&team)
		team.AddMember(guildID)
		return &team
	}
	return nil
}

func DeleteTeam(name string) {
	var team Team
	// Delete team members where teamname
	database.Where("team_name = ?", name).Delete(&TeamMember{})
	database.Where("name = ?", name).First(&team)
	database.Delete(&team)
}

func SetupTeam() {
	database.AutoMigrate(&Team{})
	database.AutoMigrate(&TeamMember{})
}

func GetTeam(name string) *Team {
	var team Team
	result := database.Where("name = ?", name).First(&team)
	if result.Error != nil {
		return nil
	}
	return &team
}

func GetTeams(guild string) []TeamMember {
	var tm []TeamMember
	database.Where("guild = ?", guild).Find(&tm)
	return tm
}

func GetTeamsLike(guild string, name string) []TeamMember {
	var tm []TeamMember
	database.Where("guild = ? AND team_name LIKE ?", guild, "%"+name+"%").Find(&tm)
	return tm
}

var (
	teamCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "team",
			Description: "Team management commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "join",
					Description: "Join a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team to join",
							Autocomplete: true,
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     false,
						},
					},
				},
				{
					Name:        "leave",
					Description: "Leave a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team to leave",
							Autocomplete: true,
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     false,
						},
					},
				},
				{
					Name:        "invite",
					Description: "Invite a guild to a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team to invite to",
							Autocomplete: true,
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     false,
						},
						{
							Name:        "guild",
							Description: "The guild to invite",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "kick",
					Description: "Kick a member from a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team to kick from",
							Type:         discordgo.ApplicationCommandOptionString,
							Autocomplete: true,
							Required:     false,
						},
						{
							Name:        "member",
							Description: "The member to kick",
							Type:        discordgo.ApplicationCommandOptionUser,
							Required:    false,
						},
					},
				},
				{
					Name:        "list",
					Description: "List all teams the guild are in",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "info",
					Description: "Get information about a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team to get information about",
							Type:         discordgo.ApplicationCommandOptionString,
							Autocomplete: true,
							Required:     true,
						},
					},
				},
				{
					Name:        "create",
					Description: "Create a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "name",
							Description: "The name of the team",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "description",
							Description: "The description of the team",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "delete",
					Description: "Delete a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team to delete",
							Type:         discordgo.ApplicationCommandOptionString,
							Autocomplete: true,
							Required:     true,
						},
					},
				},
				{
					Name:        "set",
					Description: "Set a team properties",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team to set properties for",
							Type:         discordgo.ApplicationCommandOptionString,
							Autocomplete: true,
							Required:     true,
						},
						{
							Name:        "description",
							Description: "The description of the team",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "promote",
					Description: "Promote a member to a team leader",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team where the member should promote",
							Autocomplete: true,
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     true,
						},
						{
							Name:        "member",
							Description: "The member to promote",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "demote",
					Description: "Demote a team leader to a member",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:         "team",
							Description:  "The team where the member should demote",
							Autocomplete: true,
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     true,
						},
						{
							Name:        "member",
							Description: "The member to demote",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
	}
	teamCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"team": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			rootCmd := i.ApplicationCommandData().Options[0]
			// If user don't have manage permission, cancel
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				if i.Member.Permissions&discordgo.PermissionManageServer == 0 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You don't have permission to manage teams",
						},
					})
					return
				}
				switch rootCmd.Name {
				case "create":
					description := ""
					name := rootCmd.Options[0].StringValue()
					if len(rootCmd.Options) == 2 {
						description = rootCmd.Options[1].StringValue()
					}
					team := CreateTeam(i.GuildID, name, description)
					if team == nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Failed to create team",
							},
						})
						return
					}
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "`" + name + "` created.",
						}})
					break
				case "delete":
					name := rootCmd.Options[0].StringValue()
					team := GetTeam(name)
					if team == nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Failed to delete team",
							},
						})
						return
					}
					if team.GetMember(i.GuildID).Role != TeamMemberRoleOwner {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "You must be a team owner to delete a team",
							},
						})
						return
					}
					DeleteTeam(name)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "`" + name + "` deleted.",
						},
					})
					break
				case "list":
					tms := GetTeams(i.GuildID)
					if len(tms) == 0 {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "No teams found",
							},
						})
						return
					}
					embeds := make([]*discordgo.MessageEmbed, len(tms))
					for index, tm := range tms {
						team := tm.GetTeam()
						embeds[index] = team.BuildEmbed()
					}
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds:  embeds,
							Content: "**Teams**",
						},
					})
					break
				case "set":
					name := rootCmd.Options[0].StringValue()
					team := GetTeam(name)
					if team == nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Team not found",
							},
						})
						return
					}
					if team.GetMember(i.GuildID).Role == TeamMemberRoleMember {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "You must be a team owner to delete a team",
							},
						})
						return
					}
					if len(rootCmd.Options) > 1 {
						description := rootCmd.Options[1].StringValue()
						team.Description = description
						team.Save()
					}
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Team `" + team.Name + "` changed.",
						},
					})
					break
				case "info":
					name := rootCmd.Options[0].StringValue()
					team := GetTeam(name)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "**Information**",
							Embeds:  []*discordgo.MessageEmbed{team.BuildEmbed()},
						},
					})
				}
				break
			case discordgo.InteractionApplicationCommandAutocomplete:
				var choices []*discordgo.ApplicationCommandOptionChoice
				if rootCmd.Name != "list" && rootCmd.Name != "create" && rootCmd.Options[0].Focused {
					tms := GetTeamsLike(i.GuildID, rootCmd.Options[0].StringValue())
					choices = make([]*discordgo.ApplicationCommandOptionChoice, len(tms))
					for index, tm := range tms {
						choices[index] = &discordgo.ApplicationCommandOptionChoice{
							Name:  tm.TeamName,
							Value: tm.TeamName,
						}
					}
				}
				if choices != nil {
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionApplicationCommandAutocompleteResult,
						Data: &discordgo.InteractionResponseData{
							Choices: choices,
						},
					})
					if err != nil {
						panic(err)
					}
				}
			}
		},
	}
)
