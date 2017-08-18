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

// UserIDHasRoleByGuild checks if a user by ID has a role in a guild
func UserIDHasRoleByGuild(r, u string, g *discordgo.Guild) bool {
	for _, m := range g.Members {
		if m.User.ID == u {
			rid, _ := GetRoleIDByName(r, g)
			for _, role := range m.Roles {
				if role == rid {
					return true
				}
			}
		}
	}
	return false
}

// IsDirectMessage checks if a channel is a DM
func IsDirectMessage(cid string, s *discordgo.Session) bool {
	if c, _ := s.Channel(cid); c.IsPrivate {
		return true
	}
	return false
}
