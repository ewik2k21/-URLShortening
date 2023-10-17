package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/ewik2k21/-URLShortening/cmd/config"
	"github.com/gin-gonic/gin"
)

var links = make(map[string]string)

func main() {
	config.ParseFlags()
	router := gin.Default()
	router.Use(methodSelector)

	err := router.Run(config.FlagA)
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
	if strings.Contains(config.FlagB, "http") {
		c.Writer.Write([]byte(config.FlagB + "/" + id))
		return
	}
	c.Writer.Write([]byte("http://" + config.FlagB + "/" + id))

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
