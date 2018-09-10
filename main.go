package main

import (
	"log"
	"net/http"
	"strings"
)

type jsonHandler struct {
	defaultHandler http.Handler
}

func (h *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".json") {
		w.Header().Set("Content-Type", "application/json")
	}

	h.defaultHandler.ServeHTTP(w, r)
}

func main() {
	h := &jsonHandler{http.FileServer(http.Dir("./"))}
	http.Handle("/", h)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
