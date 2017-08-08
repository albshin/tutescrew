package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Context struct to hold variables passed to command handlers
type Context struct {
	Cmd  string
	Args []string
	Msg  *discordgo.MessageCreate
	Sess *discordgo.Session
}

// Command defines the interface for commands
type Command interface {
	handle(Context) error
	description() string
	usage() string
}

// Handler holds a map of all active commands to be handled
type Handler struct {
	Commands map[string]Command
}

// AddCommand adds commands to be handled
func (h *Handler) AddCommand(c string, cmd Command) {
	h.Commands[c] = cmd
}

// Handle creates a new goroutine to handle the command
func (h *Handler) Handle(ctx Context) {
	if obj, ok := h.Commands[ctx.Cmd]; ok {
		go obj.handle(ctx)
	} else {
		ctx.Sess.ChannelMessageSend(ctx.Msg.ChannelID, "Command does not exist!")
	}
}
