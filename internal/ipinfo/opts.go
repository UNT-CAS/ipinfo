package ipinfo

import (
	"github.com/namsral/flag"
)

var (
	// Locale assists in determining the language
	Locale = flag.String("locale", "en", "locale")
	// Port to bind the http server on
	Port = flag.Int("port", 8000, "port to bind http server")
	// Loglevel (0=debug, 1=info, 2=warn, 3=error)
	Loglevel = flag.Int("loglevel", 1, "log level (0=debug, 1=info, 2=warn, 3=error)")
)
