package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/oliveagle/jsonpath"
	"github.com/pkg/browser"
)

var (
	// Version is given on compile time.
	Version string
)

// JSONHandler is a handler which returns JSON with "Content-Type: application/json".
type JSONHandler struct {
	defaultHandler http.Handler
	route          map[*regexp.Regexp]string
	redirect       map[*regexp.Regexp]string
}

// Config represents a server congig.
type Config struct {
	Port     uint16
	Route    map[string]string
	Redirect map[string]string
}

type postedJSON map[string]interface{}

func (json *postedJSON) Lookup(path string) string {
	if json == nil {
		return ""
	}

	jpath := strings.Trim(path, "{}")
	res, err := jsonpath.JsonPathLookup(*json, jpath)
	if err != nil {
		log.Println("JsonPathLookup error:", jpath, err)
		return ""
	}

	return fmt.Sprintf("%v", res)
}

var jpathPattern = regexp.MustCompile(`{\$\..+?}`)

func (h *JSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for from, to := range h.redirect {
		if from.MatchString(r.URL.Path) {
			destination := from.ReplaceAllString(r.URL.Path, to)

			http.Redirect(w, r, destination, http.StatusFound)
			return
		}
	}

	var replacer func(string) string
	switch r.Method {
	case http.MethodPost:
		mimeType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			log.Println("ParseMediaType error:", err)
		} else if mimeType == "application/json" {
			body, _ := ioutil.ReadAll(r.Body)

			var jsonData postedJSON
			if err := json.Unmarshal(body, &jsonData); err != nil {
				log.Println("Unmarshal error:", string(body), err)
			} else {
				replacer = jsonData.Lookup
			}
		}
	}
	if replacer == nil {
		replacer = func(_ string) string {
			return ""
		}
	}

	for from, to := range h.route {
		if from.MatchString(r.URL.Path) {
			file := jpathPattern.ReplaceAllStringFunc(from.ReplaceAllString(r.URL.Path, to), replacer)
			log.Println(r.URL.Path, "->", file)

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

// NewJSONHandler constructs a new JSONHandler.
func NewJSONHandler(config Config) *JSONHandler {
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

	return &JSONHandler{
		defaultHandler: http.FileServer(http.Dir("./")),
		route:          route,
		redirect:       redirect,
	}
}

// NewServer constructs a new http.Server configured with a given config.
func NewServer(config Config) *http.Server {
	var port string
	if config.Port != 0 {
		port = strconv.Itoa(int(config.Port))
	} else {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.Handle("/", NewJSONHandler(config))

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

const configJSON = "goserve.json"
const host = "http://localhost"

func main() {
	fmt.Printf("goserve %s", Version)
	fmt.Println()

	var config Config
	if raw, err := ioutil.ReadFile(configJSON); err == nil {
		if err := json.Unmarshal(raw, &config); err != nil {
			log.Println("Unmarshal error:", configJSON, err)
		}
	} else {
		log.Println(err)
	}

	server := NewServer(config)

	time.AfterFunc(500*time.Millisecond, func() {
		err := browser.OpenURL(host + server.Addr)
		if err != nil {
			log.Println(err)
		}
	})

	err := server.ListenAndServe()
	log.Fatal(err)
}
