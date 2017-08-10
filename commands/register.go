package commands

import (
	"errors"
	"fmt"
	"net/url"
)

// Register struct class
type Register struct {
	CASAuthURL     string
	CASRedirectURL string
}

func (r *Register) handle(ctx Context) error {
	// Check if student is already registered
	ch, err := ctx.Sess.State.Channel(ctx.Msg.ChannelID)
	if err != nil {
		fmt.Println("Could not get channel")
	}
	g, _ := ctx.Sess.State.Guild(ch.GuildID)

	rid, _ := getRoleIDByName("Student", g)
	mem, _ := ctx.Sess.GuildMember(ch.GuildID, ctx.Msg.Author.ID)
	for _, role := range mem.Roles {
		if role == rid {
			ctx.Sess.ChannelMessageSend(ctx.Msg.ChannelID, "You are already registered!")
			return errors.New("already registered")
		}
	}

	// Build the full login URL
	u, _ := url.Parse(r.CASAuthURL)
	q := u.Query()
	q.Add("service", r.CASRedirectURL)
	u.RawQuery = q.Encode()

	ctx.Sess.ChannelMessageSend(ctx.Msg.ChannelID, "Please go to "+u.String()+" to start the registration process.")

	return nil
}

func (r *Register) description() string {
	return "Allows the user to register for the \"student role\". A valid RCSID is required."
}
func (r *Register) usage() string { return "" }
