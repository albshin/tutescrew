package route

import (
	"github.com/albshin/tutescrew/config"
	"github.com/bwmarrin/discordgo"

	"github.com/julienschmidt/httprouter"
)

var (
	cfg     config.CASConfig
	session *discordgo.Session
)

// Router creates a new configured router
func Router(c config.CASConfig, sess *discordgo.Session) *httprouter.Router {
	r := httprouter.New()

	cfg = c
	session = sess

	r.GET("/auth/cas", handleCAS)

	return r
}
