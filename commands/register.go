package commands

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

var (
	baseURL, _ = url.Parse("http://info.rpi.edu/people/search/")
)

// Register struct class
type Register struct{}

func (r Register) handle(ctx Context) error {
	rcsid := ctx.Args[0]
	end, _ := url.Parse(rcsid)
	urlStr := baseURL.ResolveReference(end)

	resp, err := http.Get(urlStr.String())
	if err != nil {
		return errors.New("unable to retrieve info")
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return errors.New("unable to retrieve info")
	}

	// Only find first result
	// If no result is found, assume user put in rcsid wrong
	result, ok := scrape.Find(root, scrape.ByClass("search-results"))
	if !ok {
		return errors.New("invalid RCSID")
	}

	email, ok := scrape.Find(result, scrape.ByClass("field-name-field-email"))
	if !ok {
		return errors.New("could not find email")
	}
	if !strings.HasPrefix(scrape.Text(email), rcsid) {
		ctx.Sess.ChannelMessageSend(ctx.Msg.ChannelID, "Not a valid RCSID!")
		return errors.New("invalid RCSID")
	}

	status, ok := scrape.Find(result, scrape.ByClass("field-name-field-status"))
	if !ok {
		return errors.New("could not find class status")
	}
	class := scrape.Text(status)[:2]

	if class != "SR" && class != "JR" && class != "SO" && class != "FR" {
		return errors.New("listed as inactive student")
	}

	ch, err := ctx.Sess.State.Channel(ctx.Msg.ChannelID)
	if err != nil {
		fmt.Println("Could not get channel")
	}

	g, _ := ctx.Sess.State.Guild(ch.GuildID)
	var rid string
	for _, role := range g.Roles {
		if role.Name == "Student" {
			rid = role.ID
		}
	}

	ctx.Sess.GuildMemberRoleAdd(ch.GuildID, ctx.Msg.Author.ID, rid)
	ctx.Sess.ChannelMessageSend(ctx.Msg.ChannelID, "Successfully registered as a student!")
	return nil
}

func (r Register) description() string {
	return "Allows the user to register for the \"student role\". A valid RCSID is required."
}
func (r Register) usage() string { return "<rcsid>" }
