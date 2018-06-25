package clients

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/umahmood/haversine"
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

func TestGetLatLongZIPCode(t *testing.T) {
	/*zipcode := "94086"
	geoAddress := fmt.Sprintf("%s", zipcode)
	res, _ := GetLatLong(url.QueryEscape(geoAddress))
	xlog.Infof("lat :%f", res.Lat)
	xlog.Infof("lon :%f", res.Lon)*/

	oxford := haversine.Coord{Lat: 51.45, Lon: 1.15} // Oxford, UK
	turin := haversine.Coord{Lat: 45.04, Lon: 7.42}  // Turin, Italy
	mi, km := haversine.Distance(oxford, turin)
	fmt.Println("Miles:", mi, "Kilometers:", km)
}
