package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/albshin/teamRPI-bot/commands"
	"github.com/albshin/teamRPI-bot/router"
	"github.com/bwmarrin/discordgo"
)

var (
	config Configuration
	h      *commands.Handler
)

// Configuration holds config values
type Configuration struct {
	Token          string
	Prefix         string
	CASAuthURL     string
	CASRedirectURL string
}

func init() {
	// Read config file
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	log.Println("Config loaded.")
}

func main() {
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error!", err)
		return
	}

	dg.AddHandler(messageCreate)

	h = &commands.Handler{Commands: make(map[string]commands.Command)}
	h.AddCommand("register", &commands.Register{CASAuthURL: config.CASAuthURL, CASRedirectURL: config.CASRedirectURL})

	err = dg.Open()
	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("Bot is running...")

	r := router.Router(config.CASAuthURL, config.CASRedirectURL, dg)
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
	if !strings.HasPrefix(m.Content, config.Prefix) {
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
