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
	database.Where("name = ?", name).First(&team)
	database.Delete(&team)
}

func SetupTeam() {
	database.AutoMigrate(&Team{})
	database.AutoMigrate(&TeamMember{})
}

func GetTeam(name string) *Team {
	var team Team
	database.Where("name = ?", name).First(&team)
	return &team
}

func GetTeams(guild string) []TeamMember {
	var tm []TeamMember
	database.Where("guild = ?", guild).Find(&tm)
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
							Name:        "team",
							Description: "The team to join",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "leave",
					Description: "Leave a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to leave",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "invite",
					Description: "Invite a guild to a team",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to invite to",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
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
							Name:        "team",
							Description: "The team to kick from",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
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
							Name:        "team",
							Description: "The team to get information about",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
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
							Name:        "team",
							Description: "The team to delete",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "set",
					Description: "Set a team properties",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to set properties for",
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
					Name:        "promote",
					Description: "Promote a member to a team leader",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to promote",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
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
							Name:        "team",
							Description: "The team to demote",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
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
					embed := &discordgo.MessageEmbed{
						Title:       team.Name,
						Description: team.Description,
						Fields: []*discordgo.MessageEmbedField{
							{
								Name: "Members",
								// Value are team members separated with new line
								Value: strings.Join(team.GetMemberNames(), "\n"),
							},
						},
					}
					embeds[index] = embed
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds:  embeds,
						Content: "**Teams**",
					},
				})
				break
			}
		},
	}
)
