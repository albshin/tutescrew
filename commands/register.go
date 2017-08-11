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

	rid, _ := GetRoleIDByName("Student", g)
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

	// Encode Discord values into the redirect
	re, _ := url.Parse(r.CASRedirectURL)
	reque := re.Query()
	reque.Add("guild", ch.GuildID)
	reque.Add("discordID", ctx.Msg.Author.ID)
	re.RawQuery = reque.Encode()

	// Add redirect to url
	q.Add("service", re.String())
	u.RawQuery = q.Encode()

	usrch, _ := ctx.Sess.UserChannelCreate(ctx.Msg.Author.ID)
	ctx.Sess.ChannelMessageSend(usrch.ID, "Please go to "+u.String()+" to start the registration process.")

	return nil
}

func (r *Register) description() string {
	return "Allows the user to start the student validation process. Upon success, the user with receive the \"student\" role."
}
func (r *Register) usage() string { return "" }
