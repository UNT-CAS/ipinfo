package ipinfo

import (
	"github.com/namsral/flag"
)

type Opts struct {
	Locale  string
	Port    int
	Verbose bool
	Version bool
}

func ParseOpts() *Opts {
	opts := new(Opts)

	flag.StringVar(&opts.Locale, "locale", "en", "Locale")
	flag.IntVar(&opts.Port, "port", 80, "Port")
	flag.BoolVar(&opts.Verbose, "verbose", false, "Show debugging and extraneous information")
	flag.BoolVar(&opts.Version, "version", false, "Print version and exit")

	flag.Parse()

	return opts
}
