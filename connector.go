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

// Caller is a handler for calling http endpoints
type Caller interface {
	Call(URL string) ([]byte, int, error)
	CallPost(URL string, data []byte) ([]byte, int, error)
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
	token  string
	caller Caller
	sync.RWMutex
}

// NewConnector returns a new Connector instance
func NewConnector(token string, callHandler Caller) *Connector {
	return &Connector{token: token, caller: callHandler}
}

// call a GET endpoint for a given service
func (c *Connector) callGet(endpoint string, service string) (*Response, string, error) {
	url, err := parseURL(endpoint, service, c.token)
	if err != nil {
		return nil, "", err
	}

	if c.caller != nil {
		res, err := c.getCaller(url)
		return res, url, err
	}

	r, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("fetch error: %s", err)
	}
	defer r.Body.Close()

	res, err := parseResponse(r.Body)
	if err != nil {
		return nil, "", err
	}

	return res, url, nil
}

// call a GET endpoint using the given handler
func (c *Connector) getCaller(endpoint string) (*Response, error) {
	byt, status, err := c.caller.Call(endpoint)
	if err != nil {
		return &Response{}, err
	}

	if status > 299 {
		return &Response{}, fmt.Errorf("call error, status %d", status)
	}

	r := bytes.NewBuffer(byt)

	return parseResponse(r)
}

// call a POST endpoint for a given service
func (c *Connector) callPost(endpoint string, service string, data []byte) (*Response, error) {
	url, err := parseURL(endpoint, service, c.token)
	if err != nil {
		return nil, err
	}

	if c.caller != nil {
		return c.postCaller(url, data)
	}

	r, err := http.Post(url, "text/html", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("fetch error: %s", err)
	}
	defer r.Body.Close()

	return parseResponse(r.Body)
}

// call a POST endpoint using the given handler
func (c *Connector) postCaller(endpoint string, data []byte) (*Response, error) {
	byt, status, err := c.caller.CallPost(endpoint, data)
	if err != nil {
		return &Response{}, err
	}

	if status > 299 {
		return &Response{}, fmt.Errorf("call error, status %d", status)
	}

	r := bytes.NewBuffer(byt)

	return parseResponse(r)
}

func parseURL(endpoint string, service string, token string) (string, error) {
	serv, ok := services[service]
	if !ok {
		return "", fmt.Errorf("'%s' is not a valid service", service)
	}

	if len(token) < 1 {
		return "", errors.New("no token supplied")
	}

	return fmt.Sprintf("%s/%s.svc/16/0/%s/%s",
		stripSlashes(baseURL), serv, token, stripSlashes(endpoint)), nil
}

func parseResponse(body io.Reader) (*Response, error) {
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
