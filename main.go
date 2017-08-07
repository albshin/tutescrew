package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	config Configuration
)

// Configuration holds config values
type Configuration struct {
	Token  string
	Prefix string
}

func init() {
	// Read config file
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(config)
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
	if !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	splt := strings.Split(m.Content, " ")
	cmd := splt[0][1:]
	args := splt[1:]
}
