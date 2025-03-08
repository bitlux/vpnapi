package vpnapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (s Security) String() string {
	return fmt.Sprintf("VPN: %t proxy: %t Tor: %t relay: %t", s.VPN, s.Proxy, s.Tor, s.Relay)
}

// The location section of a response
type Location struct {
	City              string `json:"city"`
	Region            string `json:"region"`
	Country           string `json:"country"`
	Continent         string `json:"contient"`
	RegionCode        string `json:"region_code"`
	CountryCode       string `json:"country_code"`
	ContinentCode     string `json:"continent_code"`
	Latitude          string `json:"latitude"`
	Longitude         string `json:"longitude"`
	TimeZone          string `json:"time_zone"`
	LocaleCode        string `json:"locale_code"`
	MetroCode         string `json:"metro_code"`
	IsInEuropeanUnion bool   `json:"is_in_european_union"`
}

func (l Location) String() string {
	return fmt.Sprintf("%s, %s %s", l.City, l.Region, l.CountryCode)
}

// The network section of a response
type Network struct {
	Network                      string `json:"network"`
	AutonomousSystemNumber       string `json:"autonomous_system_number"`
	AutonomousSystemOrganization string `json:"autonomous_system_organization"`
}

func (n Network) String() string {
	return fmt.Sprintf("%s (%s)", n.AutonomousSystemOrganization, n.AutonomousSystemNumber)
}

// The response type returned from a query
type Response struct {
	IP       string   `json:"ip"`
	Security Security `json:"security"`
	Location Location `json:"location"`
	Network  Network  `json:"network"`
	Message  string   `json:"message"`
}

func (r Response) String() string {
	ret := fmt.Sprintf("IP: %s\nSecurity: %s\nLocation: %s\nNetwork: %s\n", r.IP, r.Security, r.Location, r.Network)
	if r.Message != "" {
		ret += fmt.Sprintf("Message: %s\n", r.Message)
	}
	return ret
}

// Client makes queries to the API.
type Client struct {
	apiKey  string
	verbose bool
}

// New creates a new Client. Obtain an API key by registering at vpnapi.io.
func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// SetVerbose toggles verbose output to stdout.
func (c *Client) SetVerbose(v bool) *Client {
	c.verbose = v
	return c
}

// Query queries the API for details about ip.
func (c *Client) Query(ip string) (*Response, error) {
	url := baseURL + ip + "?key="
	if c.verbose {
		fmt.Printf("Querying %s\n", url+"XXXXX")
	}
	resp, err := http.Get(url + c.apiKey)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if c.verbose {
		tokens := strings.Fields(string(body))
		tmp := resp.Status + " \"" + strings.Join(tokens, " ")
		if len(tmp) > 80 {
			tmp = tmp[0:80] + "[...]"
		}
		tmp += "\""
		fmt.Printf("Got response %s\n", tmp)
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

	if c.verbose {
		fmt.Printf("Parsed response:\n%s", ret)
	}

	return ret, nil
}
