package main

import "github.com/bwmarrin/discordgo"

type Team struct {
	Name        string `gorm:"primary_key"`
	Description string
	Members     []TeamMember
}

type TeamMember struct {
	Guild string `gorm:"primary_key"`
	Role  TeamMemberRole
}
type TeamMemberRole int

const (
	TeamMemberRoleOwner TeamMemberRole = iota
	TeamMemberRoleLeader
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

func (t *Team) GetMember(member string) *TeamMember {
	for _, m := range t.Members {
		if m.Guild == member {
			return &m
		}
	}
	return &TeamMember{}
}

func (t *TeamMember) Promote() {
	t.Role = TeamMemberRoleLeader
}
func (t *TeamMember) Demote() {
	if t.Role == TeamMemberRoleLeader {
		t.Role = TeamMemberRoleMember
	}
}

func (t *Team) Save() {
	database.Save(t)
}

func CreateTeam(name string, description string) *Team {
	team := Team{Name: name, Description: description}
	database.Create(&team)
	return &team
}

func SetupTeam() {
	database.AutoMigrate(&Team{})
}

func GetTeam(name string) *Team {
	var team Team
	database.Where("name = ?", name).First(&team)
	return &team
}

var (
	teamCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "teams",
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
		"teams": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			rootCmd := i.ApplicationCommandData().Options[0]
			switch rootCmd.StringValue() {
			case "create":
				description := ""
				name := rootCmd.Options[0].StringValue()
				if len(rootCmd.Options) == 2 {
					description = rootCmd.Options[1].StringValue()
				}
				CreateTeam(name, description)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "`" + name + "` created.",
					}})

			}
		},
	}
)
