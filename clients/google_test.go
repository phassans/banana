package clients

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGetLatLong(t *testing.T) {
	line1 := "747 Calla Dr"
	line2 := "Apt 1"
	city := "Sunnyvale"
	state := "CA"
	geoAddress := fmt.Sprintf("%s,%s,%s,%s", line1, line2, city, state)
	GetLatLong(url.QueryEscape(geoAddress))
}

func TestGetLatLongInvalid(t *testing.T) {
	line1 := "foobar"
	line2 := "Apt 1"
	city := "Sunnyvale"
	state := "CA"
	geoAddress := fmt.Sprintf("%s,%s,%s,%s", line1, line2, city, state)
	GetLatLong(url.QueryEscape(geoAddress))
}
