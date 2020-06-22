package util

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Location _location

type _location struct{}

type regeo struct {
	Status     json.Number `json:"status"`
	Info       string      `json:"info"`
	Regeocodes struct {
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Province     string          `json:"province"`
			City         json.RawMessage `json:"city"`
			CityCode     string          `json:"city_code"`
			District     string          `json:"district"`
			AdCode       string          `json:"ad_code"`
			Township     string          `json:"township"`
			Neighborhood struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"neighborhood"`
			Building struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"building"`
			StreetNumber struct {
				Street    string `json:"street"`
				Number    string `json:"number"`
				Location  string `json:"location"`
				Direction string `json:"direction"`
				Distance  string `json:"distance"`
			} `json:"street_number"`
		} `json:"address_component"`
	} `json:"regeocodes"`
}

func (*_location) FindByGeoCode(ctx *gin.Context, latitude, longitude string) (*regeo, error) {
	url := "https://restapi.amap.com/v3/geocode/regeo?key=959aef6c8cf282cf66640c5ad83d8298&location=%s,%s&extensions=base"
	resp, err := http.Get(fmt.Sprintf(url, latitude, longitude))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r regeo
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, err
}
