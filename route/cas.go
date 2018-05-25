package route

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/albshin/tutescrew/commands"
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
	log.Println("HANDLING!")

	// Get query values
	g := r.URL.Query().Get("guild")
	d := r.URL.Query().Get("discord_id")

	u, err := validateCAS(r)
	if err != nil {
		return
	}

	if u != "" {
		gld, err := session.Guild(g)
		if err != nil {
			log.Println(err)
			return
		}
		rid, err := commands.GetRoleIDByName("Student", gld)
		if err != nil {
			log.Println(err)
			return
		}
		session.GuildMemberRoleAdd(g, d, rid)

		// Remove Newcomer role
		if commands.UserIDHasRoleByGuild("Newcomer", d, gld) {
			n, err := commands.GetRoleIDByName("Newcomer", gld)
			if err != nil {
				log.Println(err)
			}
			session.GuildMemberRoleRemove(g, d, n)
		}

		w.Write([]byte("Success! You may now close this window."))
	} else {
		w.Write([]byte("Failure"))
	}
}

func validateCAS(r *http.Request) (string, error) {
	u, err := getValidationURL(r)
	if err != nil {
		return "", err
	}

	c := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	r.Header.Add("User-Agent", "TuteScrew Discord Bot")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get ticket validation")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("unable to read validation body")
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
	u, err := url.Parse(cfg.RedirectURL)
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
	u, err := url.Parse(cfg.AuthURL)
	if err != nil {
		return "", err
	}
	v, err := url.Parse("serviceValidate")
	if err != nil {
		return "", err
	}
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
