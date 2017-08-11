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
)

var (
	cfg config.Configuration
	h   *commands.Handler
)

func main() {
	err := config.LoadConfig(&cfg)
	if err != nil {
		log.Fatal("Config could not be read correctly!")
	}

	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("Error!", err)
		return
	}

	dg.AddHandler(messageCreate)

	h = &commands.Handler{Commands: make(map[string]commands.Command)}
	h.AddCommand("register", &commands.Register{Config: cfg.CASInfo})

	err = dg.Open()
	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("Bot is running...")

	r := route.Router(cfg.CASInfo, dg)
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

	// Get relevant info to pass
	/*
		ch, err := s.State.Channel(m.ChannelID)
		if err != nil {
			fmt.Println("Could not get channel")
		}
		g, err := s.State.Guild(ch.GuildID)
		if err != nil {
			fmt.Println("Could not get guild")
		}
	*/

	ctx := commands.Context{Cmd: cmd, Args: args, Msg: m, Sess: s}
	h.Handle(ctx)
}
