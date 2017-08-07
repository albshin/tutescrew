package commands

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

var (
	baseURL, _ = url.Parse("http://info.rpi.edu/people/search/")
)

// Register struct class
type Register struct{}

func handle(rcsid string) error {
	end, _ := url.Parse(rcsid)
	urlStr := baseURL.ResolveReference(end)

	resp, err := http.Get(urlStr.String())
	if err != nil {
		return errors.New("unable to retrieve info")
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return errors.New("unable to retrieve info")
	}

	// Only find first result
	// If no result is found, assume user put in rcsid wrong
	result, ok := scrape.Find(root, scrape.ByClass("search-results"))
	if !ok {
		return errors.New("invalid RCSID")
	}
	status, ok := scrape.Find(result, scrape.ByClass("field-name-field-status"))
	if !ok {
		return errors.New("invalid RCSID")
	}
	class := scrape.Text(status)[:2]

	if class != "SR" && class != "JR" && class != "SO" && class != "FR" {
		return errors.New("listed as inactive student")
	}

	return nil
}

func (r *Register) description() string {
	return "Allows the user to register for the \"student role\". A valid RCSID is required."
}
func (r *Register) usage() string { return "<rcsid>" }
