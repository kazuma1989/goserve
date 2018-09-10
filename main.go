package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/browser"
)

type jsonHandler struct {
	defaultHandler http.Handler
	routes         map[string]string
}

func (h *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.routes != nil && h.routes[r.URL.Path] != "" {
		file := h.routes[r.URL.Path]

		if strings.HasSuffix(file, ".json") {
			w.Header().Set("Content-Type", "application/json")
		}

		http.ServeFile(w, r, file)
		return
	}

	if strings.HasSuffix(r.URL.Path, ".json") {
		w.Header().Set("Content-Type", "application/json")
	}

	h.defaultHandler.ServeHTTP(w, r)
}

func main() {
	var data map[string]string

	raw, err := ioutil.ReadFile("./routes.json")
	if err != nil {
		log.Println(err)
	} else {
		json.Unmarshal(raw, &data)
	}

	h := &jsonHandler{
		defaultHandler: http.FileServer(http.Dir("./")),
		routes:         data,
	}
	http.Handle("/", h)

	browser.OpenURL("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
