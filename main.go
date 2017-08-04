package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

var (
	prefix  string
	baseURL *url.URL
)

func init() {
	prefix = "!"
	baseURL, err := url.Parse("http://info.rpi.edu/people/search/")
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error!", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("Bot is running...")

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	splt := strings.Split(m.Content, " ")
	cmd := splt[0][1:]
	content := splt[1:]

	if cmd == "register" && (len(splt) == 2) {
		handleRegister(content[0])
	}
}

func handleRegister(rcsid string, s *discordgo.Session, m *discordgo.MessageCreate) {
	end, _ := url.Parse(rcsid)
	urlStr := baseURL.ResolveReference(end)

	resp, err := http.Get(urlStr.String())
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to retrieve info. Directory may be down.")
		return
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to retrieve info. Directory may be down.")
		return
	}

	// Only find first result
	// If no result is found, assume user put in rcsid wrong
	result, ok := scrape.Find(root, scrape.ByClass("search-results"))
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "Invalid RCSID.")
		return
	}
	status, ok := scrape.Find(result, scrape.ByClass("field-name-field-status"))
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "Invalid RCSID.")
		return
	}
	class := scrape.Text(status)[:2]

	if class != "SR" || class != "JR" || class != "SO" || class != "FR" {
		s.ChannelMessageSend(m.ChannelID, "It does not seem like you are a student.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Student status confirmed!")
}
