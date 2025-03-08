package vpnapi_test

import (
	"testing"

	"github.com/bitlux/vpnapi"
)

func TestString(t *testing.T) {
	r := vpnapi.Response{
		IP: "2601:647:4800:aaaa:bbbb:0000:0000:0000",
		Security: vpnapi.Security{
			VPN:   false,
			Proxy: false,
			Tor:   false,
			Relay: false,
		},
		Location: vpnapi.Location{
			City:              "Sunnyvale",
			Region:            "California",
			Country:           "United States",
			Continent:         "North America",
			RegionCode:        "CA",
			CountryCode:       "US",
			ContinentCode:     "NA",
			Latitude:          "37.0",
			Longitude:         "-122.0",
			TimeZone:          "America/Los_Angeles",
			LocaleCode:        "en",
			MetroCode:         "807",
			IsInEuropeanUnion: false,
		},
		Network: vpnapi.Network{
			Network:                      "2601:647:4800::/48",
			AutonomousSystemNumber:       "AS7922",
			AutonomousSystemOrganization: "COMCAST-7922",
		},
		Message: "this is the message",
	}

	want := `IP: 2601:647:4800:aaaa:bbbb:0000:0000:0000
Security: VPN: false proxy: false Tor: false relay: false
Location: Sunnyvale, California US
Network: COMCAST-7922 (AS7922)
Message: this is the message
`

	if got := r.String(); got != want {
		t.Errorf("String(): %s != %s", got, want)
	}
}
