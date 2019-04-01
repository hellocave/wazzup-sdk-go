package wazzup

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var baseURL = "https://api.wazzupsoftware.com/"

var services = map[string]string{
	"activate": "ActivateService",
	"output":   "OutputService",
}

// Response represents a Wazzup API response
type Response struct {
	XMLName   xml.Name    `xml:"Result"`
	Success   bool        `xml:"IsSuccess"`
	Error     string      `xml:"ErrorMessage"`
	Contracts []*Contract `xml:"ArrayOfMediaContractSnapshot>MediaContractSnapshot"`
	Summaries []*Summary  `xml:"ArrayOfRealEstatePropertySummarySnapshot>RealEstatePropertySummarySnapshot"`
	Property  *Property   `xml:"RealEstateProperty"`
}

// Connector contains a connection to Wazzup
type Connector struct {
	token string
	sync.RWMutex
}

// NewConnector returns a new Connector instance
func NewConnector(token string) *Connector {
	return &Connector{token: token}
}

// call a GET endpoint for a given service
func (c *Connector) callGet(endpoint string, service string) (*Response, error) {
	url, err := parseURL(endpoint, service, c.token)
	if err != nil {
		return nil, err
	}

	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %s", err)
	}

	return parseResponse(r.Body)
}

// call a POST endpoint for a given service
func (c *Connector) callPost(endpoint string, service string, data []byte) (*Response, error) {
	url, err := parseURL(endpoint, service, c.token)
	if err != nil {
		return nil, err
	}

	r, err := http.Post(url, "text/html", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("fetch error: %s", err)
	}

	return parseResponse(r.Body)
}

func parseURL(endpoint string, service string, token string) (string, error) {
	serv, ok := services[service]
	if !ok {
		return "", fmt.Errorf("'%s' is not a valid service", service)
	}

	if len(token) < 1 {
		return "", errors.New("no token supplied")
	}

	return fmt.Sprintf("%s/%s.svc/1/0/%s/%s",
		stripSlashes(baseURL), serv, token, stripSlashes(endpoint)), nil
}

func parseResponse(body io.ReadCloser) (*Response, error) {
	defer body.Close()

	bts, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read error: %s", err)
	}

	var r *Response
	err = xml.Unmarshal(bts, &r)
	if err != nil {
		return nil, err
	}

	if !r.Success {
		return r, errors.New(r.Error)
	}

	return r, nil
}

// strip the leading- and trailing slash off a uri segment
func stripSlashes(s string) string {
	return strings.Trim(s, "/")
}
