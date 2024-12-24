package main

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// docker run --env-file env -v .:/usr/share/GeoIP ghcr.io/maxmind/geoipupdate
// https://github.com/maxmind/geoipupdate/blob/main/doc/docker.md
func main() {
	ip := net.ParseIP("103.116.72.17")

	countryDB, err := geoip2.Open("Country.mmdb")
	if err != nil {
		panic(err)
	}
	defer countryDB.Close()

	co, err := countryDB.Country(ip)
	if err != nil {
		panic(err)
	}

	fmt.Println("CountryDB")
	fmt.Println("country", co.Country)
	fmt.Println()

	cityDB, err := geoip2.Open("City.mmdb")
	if err != nil {
		panic(err)
	}
	defer cityDB.Close()

	ci, err := cityDB.City(ip)
	if err != nil {
		panic(err)
	}

	fmt.Println("CityDB")
	fmt.Println("country", ci.Country)
	fmt.Println("city", ci.City)
	fmt.Println()

	asnDB, err := geoip2.Open("ASN.mmdb")
	if err != nil {
		panic(err)
	}
	defer asnDB.Close()

	asn, err := asnDB.ASN(ip)
	if err != nil {
		panic(err)
	}
	fmt.Println("ASNDB")
	fmt.Println("ASN", asn)

	// output
	// CountryDB
	// country {map[de:Hongkong en:Hong Kong es:Hong Kong fr:Hong Kong ja:香港 pt-BR:Hong Kong ru:Гонконг zh-CN:香港] HK 1819730 false}
	//
	// CityDB
	// country {map[de:Hongkong en:Hong Kong es:Hong Kong fr:Hong Kong ja:香港 pt-BR:Hong Kong ru:Гонконг zh-CN:香港] HK 1819730 false}
	// city {map[en:Kowloon City zh-CN:九龙城] 1819607}
	//
	// ASNDB
	// ASN &{PRIME-SEC 400618}
}
