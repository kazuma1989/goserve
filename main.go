package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/browser"
)

var (
	// Version is given on compile time.
	Version string
)

type jsonHandler struct {
	defaultHandler http.Handler
	routes         map[*regexp.Regexp]string
}

func (h *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for from, to := range h.routes {
		if from.MatchString(r.URL.Path) {
			file := from.ReplaceAllString(r.URL.Path, to)
			log.Println(file)

			if strings.HasSuffix(file, ".json") {
				w.Header().Set("Content-Type", "application/json")
			}

			http.ServeFile(w, r, file)
			return
		}
	}

	if strings.HasSuffix(r.URL.Path, ".json") {
		w.Header().Set("Content-Type", "application/json")
	}

	h.defaultHandler.ServeHTTP(w, r)
}

func main() {
	fmt.Printf("goserve %s", Version)
	fmt.Println()

	var data map[string]string
	if raw, err := ioutil.ReadFile("./routes.json"); err == nil {
		json.Unmarshal(raw, &data)
	} else {
		log.Println(err)
	}

	routes := make(map[*regexp.Regexp]string)
	for k, v := range data {
		r, err := regexp.Compile(k)
		if err == nil {
			routes[r] = v
		}
	}

	h := &jsonHandler{
		defaultHandler: http.FileServer(http.Dir("./")),
		routes:         routes,
	}
	http.Handle("/", h)

	browserOpen := time.AfterFunc(500*time.Millisecond, func() {
		browser.OpenURL("http://localhost:8080")
	})

	err := make(chan error)
	go func() {
		err <- http.ListenAndServe(":8080", nil)
	}()

	log.Fatal(<-err)
	browserOpen.Stop()
}
