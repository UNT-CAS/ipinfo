package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jnovack/ipinfo/internal/ipinfo"
	"github.com/mattn/go-isatty"
	"github.com/namsral/flag"
	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// The GeoIP databases
var dbCity *geoip2.Reader
var dbASN *geoip2.Reader

var opts *ipinfo.Opts

func main() {

	// Initialize the database.
	var err error
	dbCity, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Error().Err(err).Msg("Error: TODO FIXME")
	}

	dbASN, err = geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		log.Error().Err(err).Msg("Error: TODO FIXME")
	}

	http.Handle("/", http.HandlerFunc(infoLookup))
	log.Info().Msg("Listening on :" + strconv.FormatInt(int64(opts.Port), 10))
	http.ListenAndServe(":"+strconv.FormatInt(int64(opts.Port), 10), nil)
}

// https://github.com/multiverse-os/ip/blob/1c436abe71f332ef3d2342c7a08a8ad25ae379b9/records.go

type codename struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ipInfo struct {
	IP           string   `json:"ip"`
	City         string   `json:"city"`
	Region       string   `json:"region"`
	Country      codename `json:"country"`
	Continent    codename `json:"continent"`
	Location     location `json:"location"`
	Postal       string   `json:"postal"`
	ASN          uint     `json:"asn"`
	Organization string   `json:"organization"`
}

func infoLookup(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	retval := "200"

	defer func() {
		// Get the current time, so that we can then calculate the execution time.
		dur := float64(float64(time.Since(start).Nanoseconds()) / 1000000)

		duration.WithLabelValues(retval).Observe(dur)
		// Log how much time it took to respond to the request, when we're done.
		if opts.Verbose == true {
			log.Printf(
				"%s %s %s %.3f",
				DefangIP(r.RemoteAddr),
				r.Method,
				r.URL.Path,
				dur,
			)
		}
	}()

	var IPAddress string
	IPAddress = strings.Split(r.URL.Path, "/")[1]

	// Set the requested IP to the user's request request IP, if we got no address.
	if IPAddress == "" || IPAddress == "self" {
		// The request is most likely being done through a reverse proxy.
		if realIP, ok := r.Header["X-Real-Ip"]; ok && len(r.Header["X-Real-Ip"]) > 0 {
			IPAddress = realIP[0]
		} else {
			// Get the real actual request IP without the trolls
			IPAddress = DefangIP(r.RemoteAddr)
		}
	}

	ip := net.ParseIP(IPAddress)
	if ip == nil {
		http.Error(w, "Invalid IP address", http.StatusBadRequest)
		retval = "400"
		return
	}

	// Query the maxmind database for that IP address.
	recCity, err := dbCity.City(ip)
	if err != nil {
		log.Error().Err(err).Msg("Error: TODO FIXME")
	}

	// Query the maxmind database for that IP address.
	recASN, err := dbASN.ASN(ip)
	if err != nil {
		log.Error().Err(err).Msg("Error: TODO FIXME")
	}

	// String containing the region/subdivision of the IP. (E.g.: Scotland, or
	// California).
	var sd string
	// If there are subdivisions for this IP, set sd as the first element in the
	// array's name.
	if recCity.Subdivisions != nil {
		sd = recCity.Subdivisions[0].Names[opts.Locale]
	}

	loc := location{
		Latitude:  recCity.Location.Latitude,
		Longitude: recCity.Location.Longitude,
	}

	country := codename{
		Code: recCity.Country.IsoCode,
		Name: recCity.Country.Names[opts.Locale],
	}

	continent := codename{
		Code: recCity.Continent.Code,
		Name: recCity.Continent.Names[opts.Locale],
	}

	// Fill up the data array with the geoip data.
	d := ipInfo{
		IP:           ip.String(),
		Country:      country,
		City:         recCity.City.Names[opts.Locale],
		Region:       sd,
		Continent:    continent,
		Postal:       recCity.Postal.Code,
		Location:     loc,
		ASN:          recASN.AutonomousSystemNumber,
		Organization: recASN.AutonomousSystemOrganization,
	}

	// Since we don't have HTML output, nor other data from geo data,
	// everything is the same if you do /8.8.8.8, /8.8.8.8/json or /8.8.8.8/geo.
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
}

// Very restrictive, but this way it shouldn't completely fuck up.
var callbackJSONP = regexp.MustCompile(`^[a-zA-Z_\$][a-zA-Z0-9_\$]*$`)

// Remove from the IP eventual [ or ], and remove the port part of the IP.
func DefangIP(ip string) string {
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)
	ss := strings.Split(ip, ":")
	ip = strings.Join(ss[:len(ss)-1], ":")
	return ip
}

func init() {
	opts = ipinfo.ParseOpts()
	if opts.Version == true {
		os.Exit(0)
	}
	ipinfo.PrometheusInit()

	if isatty.IsTerminal(os.Stdout.Fd()) {
		// Format using ConsoleWriter if running straight
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().In(time.Local)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		// Format using JSON if running as a service (or container)
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	// TODO Permit run-time setting of LogLevel
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	flag.Parse()
}
