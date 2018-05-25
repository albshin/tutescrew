package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/albshin/tutescrew/commands"
	"github.com/albshin/tutescrew/config"
	"github.com/albshin/tutescrew/route"
	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

// Configuration holds environment variables
type Configuration struct {
	Token  string `required:"true"`
	Prefix string `default:"$"`
	CAS    config.CASConfig
}

var (
	h   *commands.Handler
	cfg Configuration
)

func main() {
	// Read config from env variables
	err := envconfig.Process("tutescrew", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("Error!", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(guildMemberAdd)

	h = &commands.Handler{Commands: make(map[string]commands.Command)}
	h.AddCommand("register", &commands.Register{Config: cfg.CAS})

	err = dg.Open()
	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("Bot is running...")

	r := route.Router(cfg.CAS, dg)
	http.ListenAndServe(":8080", r)

	fmt.Println("Web server is running...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-c
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, cfg.Prefix) {
		return
	}

	msg := strings.ToLower(m.Content)
	splt := strings.Split(msg, " ")
	cmd := splt[0][1:]
	args := splt[1:]

	ctx := &commands.Context{Cmd: cmd, Args: args, Msg: m, Sess: s}
	h.Handle(*ctx)
}

func guildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	// TODO: Configurable welcome message channel and msg
	gld, _ := s.Guild(g.GuildID)
	rid, _ := commands.GetRoleIDByName("Newcomer", gld)
	s.GuildMemberRoleAdd(g.GuildID, g.User.ID, rid)

	for _, ch := range gld.Channels {
		if ch.Name == "welcome" {
			s.ChannelMessageSend(ch.ID, "Welcome "+"<@"+g.Member.User.ID+">!\n\n"+
				"If you are a student, please enter `"+cfg.Prefix+"register`"+
				" to verify your enrollment status.\nIf you are not a currently enrolled"+
				" student, please message an admin to get tagged.")
		}
	}
}
