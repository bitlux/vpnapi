package vpnapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://vpnapi.io/api/"

var ErrRateLimited = errors.New("rate limited, try again later")

// The security section of a response
type Security struct {
	VPN   bool `json:"vpn"`
	Proxy bool `json:"proxy"`
	Tor   bool `json:"tor"`
	Relay bool `json:"relay"`
}

// The location section of a response
type Location struct {
	City              string `json:"city"`
	Region            string `json:"region"`
	Country           string `json:"country"`
	Continent         string `json:"contient"`
	RegionCode        string `json:"region_code"`
	ContinentCode     string `json:"continent_code"`
	Latitude          string `json:"latitude"`
	Longitude         string `json:"longitude"`
	TimeZone          string `json:"time_zone"`
	LocaleCode        string `json:"locale_code"`
	MetroCode         string `json:"metro_code"`
	IsInEuropeanUnion bool   `json:"is_in_european_union"`
}

// The network section of a response
type Network struct {
	Network                      string `json:"network"`
	AutonomousSystemNumber       string `json:"autonomous_system_number"`
	AutonomousSystemOrganization string `json:"autonomous_system_organization"`
}

// The response type returned from a query
type Response struct {
	IP       string   `json:"ip"`
	Security Security `json:"security"`
	Location Location `json:"location"`
	Network  Network  `json:"network"`
}

// Client makes queries to the API.
type Client struct {
	apiKey string
}

// New creates a new Client. Obtain an API key by registering at vpnapi.io.
func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// Query queries the API for details about ip.
func (c *Client) Query(ip string) (*Response, error) {
	resp, err := http.Get(baseURL + ip + "?" + c.apiKey)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode; code > 299 {
		if code == 429 {
			return nil, ErrRateLimited
		}
		return nil, fmt.Errorf("status code %d", code)
	}

	ret := &Response{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}

	return ret, nil
}
