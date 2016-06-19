package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	db   = flag.String("redis", "127.0.0.1:6379", "")
	port = flag.String("port", "8301", "")
)

func main() {
	flag.Parse()

	s := NewShort(*db)
	defer s.Close()

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		w.Header().Set("Content-Type", "application/json")

		r.ParseForm()
		url, ok := r.Form["url"]
		if !ok || len(url) == 0 || url[0] == "" {
			fmt.Fprint(w, "{\"val\":\"\", \"err\":\"invalid data\"}")
			return
		}

		val, err := s.Encode(url[0])
		if err != nil {
			fmt.Fprint(w, "{\"val\":\"\", \"err\":\"encode failed\"}")
		} else {
			fmt.Fprintf(w, "{\"val\":\"%s\", \"err\":\"\"}", val)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}

		val := r.URL.Path[1:]
		strings.TrimSpace(val)
		if val == "" {
			f, err := os.Open("asserts/index.html")
			if err != nil {
				return
			}
			defer f.Close()
			io.Copy(w, f)
			return
		}

		url, err := s.Decode(val)
		if err != nil {
			url = r.URL.Host
		}
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:"+*port, nil))
}
