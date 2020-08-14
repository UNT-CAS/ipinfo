package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jnovack/ipinfo/internal/ipinfo"
	"github.com/mattn/go-isatty"
	"github.com/namsral/flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
		// TODO Add Favicon Functionality
		// bytes, err := base64.StdEncoding.DecodeString(favicon)
		// if err != nil {
		// 	log.Error().Err(err).Msg("Unable to decode 'favicon' variable")
		// }
		// w.Header().Set("Content-Type", "image/png")
		// w.Write(bytes)
	})
	http.Handle("/", http.HandlerFunc(ipinfo.Lookup))
	log.Info().Msg("Listening on :" + strconv.FormatInt(int64(*ipinfo.Port), 10))
	http.ListenAndServe(":"+strconv.FormatInt(int64(*ipinfo.Port), 10), nil)
}

func init() {
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

	flag.Parse()

	var zerologlevel zerolog.Level
	switch *ipinfo.Loglevel {
	case -1:
		zerologlevel = zerolog.TraceLevel
	case 0:
		zerologlevel = zerolog.DebugLevel
	case 2:
		zerologlevel = zerolog.WarnLevel
	case 3:
		zerologlevel = zerolog.ErrorLevel
	default:
		zerologlevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(zerologlevel)

}

var favicon = "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABUklEQVQ4T63TTUsCYRDA8f9q5rYL1kEiEOwcLRWISKeI+gpF+BH8AL3hKQi8SZfo0KXPIC146hiCUG7QoYUlqG5FSJqr7kusUrGSZeWcZ37MMzOP4Lquyz9CGCpgGAa5XI58Po8sywP15eugXC6TyWRQVZVoNPp7YKCKniRfB7quk06nKZVKHNy0qbah8mxzV3dQJoLszYWJywEf4QM0TSORSGCaJjtam9MHq1M0Mx7kSG9xXbVRlyVGA8IH8i1Qs+AwKXaSa5ZLqljnOCWyGB0ZDAgFYH++C3jXkizWyC2IrE79ADQaDXavLEpPNoUliUhIoHBvsV0xOVuRmBQ/5/DlE96B80ebpgOxMQH9xSGrhNmYDvUfojc8bxOKorB12SQSgvX4KLevDrORADHJvwFP6nvKmxdmB8gq3Rn0i77AidFCCgqs9bTcCw33M/3llN8Aw17C0TKxYeUAAAAASUVORK5CYII="
