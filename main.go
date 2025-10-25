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

var Version = "dev"

var (
	addr    string        = "0.0.0.0:8000"
	refresh time.Duration = 100 * time.Millisecond
)

var (
	clearTerm     = []byte("\x1b[2J")
	cursorTopLeft = []byte("\x1b[1;1H")
	hideCursor    = []byte("\x1b[?25l")
)

const maxInt = int(^uint(0) >> 1)

func joinBytes(s [][]byte) []byte {
	if len(s) == 0 {
		return []byte{}
	} else if len(s) == 1 {
		return append([]byte(nil), s[0]...)
	}

	var n int
	for _, v := range s {
		if len(v) > maxInt-n {
			panic("bytes: Join output length overflow")
		}
		n += len(v)
	}

	b := make([]byte, n)
	bp := copy(b, s[0])
	for _, v := range s[1:] {
		bp += copy(b[bp:], v)
	}
	return b
}

func randFg() []byte {
	fg := (rand.Int() % 5) + 90
	return []byte("\x1b[" + fmt.Sprintf("%d", fg) + "m")
}

func randBg() []byte {
	bg := (rand.Int() % 6) + 41
	return []byte("\x1b[" + fmt.Sprintf("%d", bg) + "m")
}

func duck(x int, y int) []byte {
	d := fmt.Appendf(nil, `%[1]s%[2]s    _
%[2]s __(.)<
%[2]s \_)_)
`, strings.Repeat("\n", y), strings.Repeat(" ", x))
	return joinBytes([][]byte{clearTerm, cursorTopLeft, d})
}

func streamData(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	userAgent := r.Header.Get("user-agent")
	if !strings.HasPrefix(userAgent, "curl") {
		log.Printf("no curl, no 🦆 %s", userAgent)
		http.Error(w, "curl me", http.StatusBadRequest)
		return
	}

	log.Printf("new 🦆 connection from %s %s", r.RemoteAddr, userAgent)

	params := r.URL.Query()
	enableRandBg := params.Get("bg") == "1"
	enableRandFg := params.Get("fg") == "1"

	done := make(chan struct{})

	go func() {
		defer close(done)

		_, err := w.Write(joinBytes([][]byte{hideCursor, clearTerm, cursorTopLeft}))
		if err != nil {
			return
		}

		x := 0
		y := 0
		for {
			select {
			case <-done:
				return
			default:
				if x < 64 {
					x += 1
				} else if y < 12 {
					x = 0
					y += 1
				} else {
					x = 0
					y = 0
				}

				d := duck(x, y)

				if enableRandBg {
					d = joinBytes([][]byte{d, randBg()})
				}
				if enableRandFg {
					d = joinBytes([][]byte{d, randFg()})
				}

				_, err := w.Write(d)
				if err != nil {
					return
				}

				flusher.Flush()
				time.Sleep(refresh)
			}
		}
	}()

	<-r.Context().Done()
	done <- struct{}{}

	log.Printf("closed 🦆 connection from %s", r.RemoteAddr)
}

func startServer() {
	log.Printf("starting 🦆 %s on %s", Version, addr)

	http.HandleFunc("/", streamData)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	flag.StringVar(&addr, "addr", addr, "http address")
	flag.DurationVar(&refresh, "refresh", refresh, "refresh interval")

	flag.Parse()

	startServer()
}
