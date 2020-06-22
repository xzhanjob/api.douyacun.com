package util

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Location _location

type _location struct{}

type amapString string

func (a *amapString) UnmarshalJSON(b []byte) error {
	if string(b) == "[]" {
		*a = ""
	} else {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*a = amapString(s)
	}
	return nil
}

type addressComponent struct {
	Province     string     `json:"province"`
	City         amapString `json:"city"`
	CityCode     string     `json:"citycode"`
	District     string     `json:"district"`
	AdCode       string     `json:"adcode"`
	Towncode     string     `json:"towncode"`
	Neighborhood struct {
		Name amapString `json:"name"`
		Type amapString `json:"type"`
	} `json:"neighborhood"`
	Building struct {
		Name amapString `json:"name"`
		Type amapString `json:"type"`
	} `json:"building"`
	StreetNumber struct {
		Street    string `json:"street"`
		Number    string `json:"number"`
		Location  string `json:"location"`
		Direction string `json:"direction"`
		Distance  string `json:"distance"`
	} `json:"streetNumber"`
}

type regeo struct {
	Status    json.Number `json:"status"`
	Info      string      `json:"info"`
	Regeocode struct {
		FormattedAddress string           `json:"formatted_address"`
		AddressComponent addressComponent `json:"addressComponent"`
	} `json:"regeocode"`
}

func (*_location) FindByGeoCode(ctx *gin.Context, latitude, longitude string) (*addressComponent, error) {
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
	address := &r.Regeocode.AddressComponent

	return address, err
}
