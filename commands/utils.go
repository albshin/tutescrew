package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// GetRoleIDByName returns a role's role ID.
func GetRoleIDByName(n string, g *discordgo.Guild) (string, error) {
	for _, role := range g.Roles {
		if role.Name == n {
			return role.ID, nil
		}
	}
	return "", errors.New("cannot find specified role")
}

// UserHasRoleByGuild checks if a user has a role in a guild
func UserHasRoleByGuild(r, u string, g *discordgo.Guild) bool {
	for _, role := range g.Roles {
		if role.Name == r {
			return true
		}
	}
	return false
}
