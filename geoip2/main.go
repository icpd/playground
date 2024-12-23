package main

import (
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	// https://github.com/Dreamacro/maxmind-geoip/releases/latest/download/Country.mmdb
	db, err := geoip2.Open("Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP("38.9.147.134")
	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("country", record.Country)
}
