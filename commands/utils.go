package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// GetRoleIDByName returns a role's role ID.
func GetRoleIDByName(roleName string, g *discordgo.Guild) (string, error) {
	for _, role := range g.Roles {
		if role.Name == roleName {
			return role.ID, nil
		}
	}
	return "", errors.New("cannot find specified role")
}

// UserIDHasRoleByGuild checks if a user by ID has a role in a guild
func UserIDHasRoleByGuild(roleName, userID string, g *discordgo.Guild) bool {
	for _, m := range g.Members {
		if m.User.ID == userID {
			rid, _ := GetRoleIDByName(roleName, g)
			for _, r := range m.Roles {
				if r == rid {
					return true
				}
			}
		}
	}
	return false
}

// IsDirectMessage checks if a channel is a DM
func IsDirectMessage(channelID string, s *discordgo.Session) bool {
	if c, _ := s.Channel(channelID); c.IsPrivate {
		return true
	}
	return false
}
