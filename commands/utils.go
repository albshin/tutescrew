package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func getRoleIDByName(name string, g *discordgo.Guild) (string, error) {
	for _, role := range g.Roles {
		if role.Name == name {
			return role.ID, nil
		}
	}
	return "", errors.New("cannot find specified role")
}
