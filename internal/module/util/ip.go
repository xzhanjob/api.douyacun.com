package util

import (
	"dyc/internal/config"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/oschwald/geoip2-golang"
	"log"
)

var (
	ipipdb  *ipdb.City
	geoipdb *geoip2.Reader
)

//func init() {
//	ipipdb, _ = ipdb.NewCity(config.GetKey("ip::ipip_file").String())
//	geoipdb, _ = geoip2.Open(config.GetKey("ip::geo_file").String())
//}

func IPIP(ip string) (map[string]string, error) {
	return ipipdb.FindMap(ip, "CN")
}

func GeoIP(ip string) (map[string]string, error) {
	log.Printf("%s", config.GetKey("ip::geo_file").String())
	//record, err := geoipdb.City(net.ParseIP(ip))
	//if err != nil {
	//	return nil, err
	//}
	//_ = record.City.Names
	res := make(map[string]string)

	return res, nil
}
