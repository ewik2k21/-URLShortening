package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var links = make(map[string]string)

func main() {
	router := gin.Default()
	router.Use(methodSelector)
	err := router.Run(`:8080`)
	if err != nil {
		log.Fatal(err)
	}
}
func methodSelector(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		postURL(c)
	case http.MethodGet:
		getURL(c)
	}
}

func getURL(c *gin.Context) {
	id := strings.TrimPrefix(c.Request.URL.Path, "/")
	c.Writer.Header().Set("Location", links[id])
	c.Status(http.StatusTemporaryRedirect)
}

func postURL(c *gin.Context) {
	id := GenerateUniqeString(8)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(400)
		return
	}
	links[id] = string(body)
	c.Status(http.StatusCreated)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Write([]byte("http://localhost:8080/" + id))

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
