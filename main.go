package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
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
	route          map[*regexp.Regexp]string
	redirect       map[*regexp.Regexp]string
}

type goserveConfig struct {
	Port     uint16
	Route    map[string]string
	Redirect map[string]string
}

func (h *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for from, to := range h.redirect {
		if from.MatchString(r.URL.Path) {
			destination := from.ReplaceAllString(r.URL.Path, to)

			http.Redirect(w, r, destination, http.StatusFound)
			return
		}
	}

	for from, to := range h.route {
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

	var config goserveConfig
	if raw, err := ioutil.ReadFile("./goserve.json"); err == nil {
		json.Unmarshal(raw, &config)
	} else {
		log.Println(err)
	}

	route := make(map[*regexp.Regexp]string)
	for k, v := range config.Route {
		r, err := regexp.Compile(k)
		if err == nil {
			route[r] = v
		}
	}

	redirect := make(map[*regexp.Regexp]string)
	for k, v := range config.Redirect {
		r, err := regexp.Compile(k)
		if err == nil {
			redirect[r] = v
		}
	}

	handler := &jsonHandler{
		defaultHandler: http.FileServer(http.Dir("./")),
		route:          route,
		redirect:       redirect,
	}
	http.Handle("/", handler)

	var port string
	if config.Port != 0 {
		port = strconv.Itoa(int(config.Port))
	} else {
		port = "8080"
	}
	browserOpen := time.AfterFunc(500*time.Millisecond, func() {
		browser.OpenURL("http://localhost:" + port)
	})

	err := make(chan error)
	go func() {
		err <- http.ListenAndServe(":"+port, nil)
	}()

	log.Fatal(<-err)
	browserOpen.Stop()
}
