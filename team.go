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

func (t *Team) GetMember(member string) TeamMember {
	for _, m := range t.Members {
		if m.Guild == member {
			return m
		}
	}
	return TeamMember{}
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
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to join",
						},
					},
				},
				{
					Name:        "leave",
					Description: "Leave a team",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to leave",
						},
					},
				},
				{
					Name:        "invite",
					Description: "Invite a guild to a team",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to invite to",
						},
						{
							Name:        "guild",
							Description: "The guild to invite",
						},
					},
				},
				{
					Name:        "kick",
					Description: "Kick a member from a team",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to kick from",
						},
						{
							Name:        "member",
							Description: "The member to kick",
						},
					},
				},
				{
					Name:        "list",
					Description: "List all teams the guild are in",
				},
				{
					Name:        "info",
					Description: "Get information about a team",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to get information about",
						},
					},
				},
				{
					Name:        "create",
					Description: "Create a team",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "name",
							Description: "The name of the team",
						},
					},
				},
				{
					Name:        "delete",
					Description: "Delete a team",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to delete",
						},
					},
				},
				{
					Name:        "set",
					Description: "Set a team properties",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to set properties for",
						},
						{
							Name:        "description",
							Description: "The description of the team",
						},
					},
				},
				{
					Name:        "promote",
					Description: "Promote a member to a team leader",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to promote",
						},
						{
							Name:        "member",
							Description: "The member to promote",
						},
					},
				},
				{
					Name:        "demote",
					Description: "Demote a team leader to a member",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "team",
							Description: "The team to demote",
						},
						{
							Name:        "member",
							Description: "The member to demote",
						},
					},
				},
			},
		},
	}
)
