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
	routes         map[*regexp.Regexp]string
}

type goserveConfig struct {
	Port   uint16
	Routes map[string]string
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

	var config goserveConfig
	if raw, err := ioutil.ReadFile("./goserve.json"); err == nil {
		json.Unmarshal(raw, &config)
	} else {
		log.Println(err)
	}

	routes := make(map[*regexp.Regexp]string)
	for k, v := range config.Routes {
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
