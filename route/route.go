package route

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/albshin/teamRPI-bot/commands"
	"github.com/bwmarrin/discordgo"

	"github.com/julienschmidt/httprouter"
)

var (
	casAuthURL     string // URL to authenticate to
	casRedirectURL string // URL to redirect to
	session        *discordgo.Session
)

// Refer to https://github.com/go-cas/cas/blob/fd85b5ae8a14a3e4a3989833499caecb4542fdf0/xml_service_response.go
// if any more struct info is needed

type casServiceResponse struct {
	XMLName xml.Name `xml:"http://www.yale.edu/tp/cas serviceResponse"`
	Failure *casAuthenticationFailure
	Success *casAuthenticationSuccess
}

type casAuthenticationFailure struct {
	XMLName xml.Name `xml:"authenticationFailure"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",innerxml"`
}

type casAuthenticationSuccess struct {
	XMLName xml.Name `xml:"authenticationSuccess"`
	User    string   `xml:"user"`
}

// Router creates a new configured router
func Router(auth, redirect string, sess *discordgo.Session) *httprouter.Router {
	r := httprouter.New()

	casAuthURL = auth
	casRedirectURL = redirect
	session = sess

	r.GET("/auth/cas", handleCAS)
	r.GET("/auth/cas/success", handleCASSuccess)

	return r
}

func handleCAS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get query values
	gid := r.URL.Query().Get("guild")
	usrid := r.URL.Query().Get("discordID")

	//TODO: Check if user has Student role

	// Validate received ticket
	u, _ := url.Parse(casAuthURL)
	val, _ := url.Parse("serviceValidate")
	u = u.ResolveReference(val)

	q := u.Query()
	srvc, _ := url.Parse(casRedirectURL)
	srvcq := srvc.Query()
	srvcq.Add("discordID", usrid)
	srvcq.Add("guild", gid)
	srvc.RawQuery = srvcq.Encode()
	q.Add("ticket", r.FormValue("ticket"))
	q.Add("service", srvc.String())
	u.RawQuery = q.Encode()

	c := &http.Client{
		Timeout: time.Second * 10,
	}
	req, _ := http.NewRequest("GET", u.String(), nil)
	r.Header.Add("User-Agent", "teamRPI Discord Bot")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Failed to get ticket validation")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	uname := strings.ToLower(handleCASResponse(body))
	if uname != "" {
		g, _ := session.Guild(gid)
		rid, _ := commands.GetRoleIDByName("Student", g)
		session.GuildMemberRoleAdd(gid, usrid, rid)
		w.Write([]byte("Success! You may now close this window."))
	} else {
		w.Write([]byte("Failure"))
	}
	// TODO: store username in database

	//http.Redirect(w, r, "/auth/cas/success", http.StatusTemporaryRedirect)
}

func handleCASSuccess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Success"))
}

// Improve error handling
func handleCASResponse(r []byte) string {
	var c casServiceResponse

	err := xml.Unmarshal(r, &c)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%+v\n", c.Failure)

	if c.Failure != nil {
		// Log
		return ""
	}

	return c.Success.User
}