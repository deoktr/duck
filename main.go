package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const VERSION = "1.0.0"

var (
	addr    string
	refresh time.Duration
)

var (
	CLEAR           = []byte("\x1b[2J")
	CURSOR_TOP_LEFT = []byte("\x1b[1;1H")
	HIDE_CURSOR     = []byte("\x1b[?25l")
)

func randFg() []byte {
	fg := (rand.Int() % 5) + 90
	return []byte("\x1b[" + fmt.Sprintf("%d", fg) + "m")
}

func randBg() []byte {
	bg := (rand.Int() % 6) + 41
	return []byte("\x1b[" + fmt.Sprintf("%d", bg) + "m")
}

func termHandler(w http.ResponseWriter, r *http.Request) {
	// black fore and background
	w.Header().Add("Fart", "\x1b[30m \x1b[40m")

	// w.Header().Add("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// hide cursor
	_, _ = w.Write(HIDE_CURSOR)
	_, _ = w.Write(CLEAR)
	_, _ = w.Write(CURSOR_TOP_LEFT)

	nl := 0
	padding := 0
	for {
		if padding < 64 {
			padding += 1
		} else if nl > 12 {
			padding = 0
			nl = 0
		} else {
			padding = 0
			nl += 1
		}

		_, _ = w.Write(CLEAR)
		_, _ = w.Write(CURSOR_TOP_LEFT)

		_, _ = w.Write([]byte("\n" + fmt.Sprintf(`%[1]s%[2]s    _
%[2]s __(.)<
%[2]s \_)_)
`, strings.Repeat("\n", nl), strings.Repeat(" ", padding))))

		// random background and foreground
		_, _ = w.Write(randFg())
		_, _ = w.Write(randBg())

		// show output in terminal
		flusher.Flush()
		time.Sleep(refresh)
	}
}

func startServer() {
	log.Printf("starting Duck 🦆 %s on %s", VERSION, addr)

	http.HandleFunc("/", termHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	flag.StringVar(&addr, "addr", "localhost:8000", "http service address")
	flag.DurationVar(&refresh, "refresh", 100*time.Millisecond, "refresh interval")
	flag.Parse()

	startServer()
}
