package main

import (
	"encoding/json"
	"fmt"
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
)

// The GeoIP databases
var dbCity *geoip2.Reader
var dbASN *geoip2.Reader

func main() {
	// Initialize the database.
	var err error
	dbCity, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	dbASN, err = geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	port := flag.Int("port", 80, "Port")
	flag.Parse()

	// http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", http.HandlerFunc(infoLookup))
	log.Info("Serving metrics on " + strconv.FormatInt(int64(*port), 10))
	log.Fatal(http.ListenAndServe(":"+strconv.FormatInt(int64(*port), 10), nil))

}

var invalidIPBytes = []byte("Please provide a valid IP address.")

type dataStruct struct {
	IP            string `json:"ip"`
	City          string `json:"city"`
	Region        string `json:"region"`
	Country       string `json:"country"`
	CountryFull   string `json:"country_full"`
	Continent     string `json:"continent"`
	ContinentFull string `json:"continent_full"`
	Loc           string `json:"loc"`
	Postal        string `json:"postal"`
}

var nameToField = map[string]func(dataStruct) string{
	"ip":             func(d dataStruct) string { return d.IP },
	"city":           func(d dataStruct) string { return d.City },
	"region":         func(d dataStruct) string { return d.Region },
	"country":        func(d dataStruct) string { return d.Country },
	"country_full":   func(d dataStruct) string { return d.CountryFull },
	"continent":      func(d dataStruct) string { return d.Continent },
	"continent_full": func(d dataStruct) string { return d.ContinentFull },
	"loc":            func(d dataStruct) string { return d.Loc },
	"postal":         func(d dataStruct) string { return d.Postal },
}

func infoLookup(w http.ResponseWriter, r *http.Request) {
	// Get the current time, so that we can then calculate the execution time.
	start := time.Now()

	// Log how much time it took to respond to the request, when we're done.
	defer log.Printf(
		"[rq] %s %s %dns",
		r.Method,
		r.URL.Path,
		time.Since(start).Nanoseconds())

	// Separate two strings when there is a / in the URL requested.
	requestedThings := strings.Split(r.URL.Path, "/")

	var IPAddress string
	var Which string
	switch len(requestedThings) {
	case 3:
		Which = requestedThings[2]
		fallthrough
	case 2:
		IPAddress = requestedThings[1]
	}

	// Set the requested IP to the user's request request IP, if we got no address.
	if IPAddress == "" || IPAddress == "self" {
		// The request is most likely being done through a reverse proxy.
		if realIP, ok := r.Header["X-Real-Ip"]; ok && len(r.Header["X-Real-Ip"]) > 0 {
			IPAddress = realIP[0]
		} else {
			// Get the real actual request IP without the trolls
			IPAddress = UnfuckRequestIP(r.RemoteAddr)
		}
	}

	ip := net.ParseIP(IPAddress)
	if ip == nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(invalidIPBytes)
		return
	}

	// Query the maxmind database for that IP address.
	record, err := dbCity.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	// String containing the region/subdivision of the IP. (E.g.: Scotland, or
	// California).
	var sd string
	// If there are subdivisions for this IP, set sd as the first element in the
	// array's name.
	if record.Subdivisions != nil {
		sd = record.Subdivisions[0].Names["en"]
	}

	// Fill up the data array with the geoip data.
	d := dataStruct{
		IP:            ip.String(),
		Country:       record.Country.IsoCode,
		CountryFull:   record.Country.Names["en"],
		City:          record.City.Names["en"],
		Region:        sd,
		Continent:     record.Continent.Code,
		ContinentFull: record.Continent.Names["en"],
		Postal:        record.Postal.Code,
		Loc:           fmt.Sprintf("%.4f,%.4f", record.Location.Latitude, record.Location.Longitude),
	}

	// Since we don't have HTML output, nor other data from geo data,
	// everything is the same if you do /8.8.8.8, /8.8.8.8/json or /8.8.8.8/geo.
	if Which == "" || Which == "json" || Which == "geo" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		callback := r.URL.Query().Get("callback")
		enableJSONP := callback != "" && len(callback) < 2000 && callbackJSONP.MatchString(callback)
		if enableJSONP {
			_, err = w.Write([]byte("/**/ typeof " + callback + " === 'function' " +
				"&& " + callback + "("))
			if err != nil {
				return
			}
		}
		enc := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "1" {
			enc.SetIndent("", "  ")
		}
		enc.Encode(d)
		if enableJSONP {
			w.Write([]byte(");"))
		}
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if val := nameToField[Which]; val != nil {
			w.Write([]byte(val(d)))
		} else {
			w.Write([]byte("undefined"))
		}
	}
}

// Very restrictive, but this way it shouldn't completely fuck up.
var callbackJSONP = regexp.MustCompile(`^[a-zA-Z_\$][a-zA-Z0-9_\$]*$`)

// Remove from the IP eventual [ or ], and remove the port part of the IP.
func UnfuckRequestIP(ip string) string {
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)
	ss := strings.Split(ip, ":")
	ip = strings.Join(ss[:len(ss)-1], ":")
	return ip
}
