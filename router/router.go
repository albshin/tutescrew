package router

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	casAuthURL     string // URL to authenticate to
	casRedirectURL string // URL to redirect to
)

// Router creates a new configured router
func Router(auth, redirect string) *httprouter.Router {
	r := httprouter.New()

	casAuthURL = auth
	casRedirectURL = redirect

	r.GET("/auth/cas", handleCAS)
	r.GET("/auth/cas/success", handleCASSuccess)

	return r
}

func handleCAS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Validate received ticket
	u, _ := url.Parse(casAuthURL)
	val, _ := url.Parse("serviceValidate")
	u = u.ResolveReference(val)

	q := u.Query()
	srvc, _ := url.Parse(casRedirectURL)
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

	fmt.Println(string(body))

	http.Redirect(w, r, "/auth/cas/success", http.StatusTemporaryRedirect)
}

func handleCASSuccess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Success"))
}
