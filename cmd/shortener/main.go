package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var links = make(map[string]string)

func main() {
	http.HandleFunc("/", methodSelector)
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}

func methodSelector(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postURL(w, r)
	case http.MethodGet:
		getURL(w, r)
	}
}

func postURL(w http.ResponseWriter, r *http.Request) {
	id := GenerateUniqeString(8)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	links[id] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + id))
}

func getURL(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.RequestURI, "/")
	w.Header().Set("Location", links[id])
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// func for generate string (id) for Get method get/{id}
func GenerateUniqeString(lenght int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	var b strings.Builder
	for i := 0; i < lenght; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
