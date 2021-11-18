package main

import "github.com/bwmarrin/discordgo"

type Team struct {
	Name        string `gorm:"primary_key"`
	Description string
	Members     []string
}

func (t *Team) AddMember(member string) {
	t.Members = append(t.Members, member)
}

func (t *Team) RemoveMember(member string) {
	for i, m := range t.Members {
		if m == member {
			t.Members = append(t.Members[:i], t.Members[i+1:]...)
			break
		}
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
			Name:        "team",
			Description: "Team management commands",
		},
	}
)
