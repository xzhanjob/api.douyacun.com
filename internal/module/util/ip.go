package util

import (
	"dyc/internal/config"
	"fmt"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/oschwald/geoip2-golang"
	"net"
)

var (
	ipipdb  *ipdb.City
	geoipdb *geoip2.Reader
)

const (
	Language = "zh-CN"
)

func Init() {
	ipipdb, _ = ipdb.NewCity(config.GetKey("ip::ipip_file").String())
	geoipdb, _ = geoip2.Open(config.GetKey("ip::geo_file").String())
}

func IPIP(ip string) (map[string]string, error) {
	  info, err := ipipdb.FindMap(ip, "CN")
	  ipipdb.FindInfo()
	  if err != nil {
	  	return nil, err
	  }
	  res := make(map[string]string)
	  res["city"] = info["city_name"]
	  res["country"] = info["country_name"]
	  return res, nil
}

func GeoIP(ip string) (map[string]string, error) {
	record, err := geoipdb.City(net.ParseIP(ip))
	if err != nil {
		return nil, err
	}
	_ = record.City.Names
	res := make(map[string]string)
	res["city"] = record.City.Names[Language]
	res["continent"] = record.Continent.Names[Language]
	res["country"] = record.Country.Names[Language]
	res["latitude"] = fmt.Sprintf("%f", record.Location.Latitude)
	res["longitude"] = fmt.Sprintf("%f", record.Location.Latitude)
	return res, nil
}
