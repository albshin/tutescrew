package route

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/albshin/tutescrew/commands"
	"github.com/albshin/tutescrew/database"
	"github.com/julienschmidt/httprouter"
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

func handleCAS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get query values
	g := r.URL.Query().Get("guild")
	d := r.URL.Query().Get("discord_id")

	u, err := validateCAS(r)
	if err != nil {
		return
	}

	//Check if RCSID is mapped to a different account
	if database.IsRegistered(u, d) {
		w.Write([]byte("This RCSID is already registered to an account!"))
		return
	}

	if u != "" {
		gld, _ := session.Guild(g)
		rid, _ := commands.GetRoleIDByName("Student", gld)
		session.GuildMemberRoleAdd(g, d, rid)

		// Remove Newcomer role
		if commands.UserIDHasRoleByGuild("Newcomer", d, gld) {
			n, _ := commands.GetRoleIDByName("Newcomer", gld)
			session.GuildMemberRoleRemove(g, d, n)
		}

		w.Write([]byte("Success! You may now close this window."))
	} else {
		w.Write([]byte("Failure"))
	}

	database.AddStudent(u, d)
}

func validateCAS(r *http.Request) (string, error) {
	u, _ := getValidationURL(r)

	c := &http.Client{
		Timeout: time.Second * 10,
	}
	req, _ := http.NewRequest("GET", u, nil)
	r.Header.Add("User-Agent", "TuteScrew Discord Bot")

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

	return strings.ToLower(handleCASResponse(body)), nil
}

func handleCASResponse(r []byte) string {
	var c casServiceResponse

	err := xml.Unmarshal(r, &c)
	if err != nil {
		log.Fatal(err)
	}

	if c.Failure != nil {
		// Log
		return ""
	}

	return c.Success.User
}
func getServiceURL(r *http.Request) (string, error) {
	u, err := url.Parse(cfg.CASRedirectURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("discord_id", r.URL.Query().Get("discord_id"))
	q.Add("guild", r.URL.Query().Get("guild"))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func getValidationURL(r *http.Request) (string, error) {
	u, _ := url.Parse(cfg.CASAuthURL)
	v, _ := url.Parse("serviceValidate")
	u = u.ResolveReference(v)

	srv, err := getServiceURL(r)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("ticket", r.FormValue("ticket"))
	q.Add("service", srv)
	u.RawQuery = q.Encode()

	return u.String(), nil
}
