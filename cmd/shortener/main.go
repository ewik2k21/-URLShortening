package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var links = make(map[string]string)

func main() {
	http.HandleFunc("/", postURL)
	http.HandleFunc("/get/", getURL)
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}

func postURL(w http.ResponseWriter, r *http.Request) {
	id := GenerateUniqeString()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	links[id] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + links[id]))
}
func getURL(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/get/")
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(links[id]))
}

// func for generate string (id) for Get method get/{id}
func GenerateUniqeString() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	lenght := 8
	var b strings.Builder
	for i := 0; i < lenght; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
