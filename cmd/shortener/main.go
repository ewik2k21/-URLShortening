package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/ewik2k21/-URLShortening/cmd/config"
	"github.com/gin-gonic/gin"
)

type Links struct {
	mu    sync.Mutex
	links map[string]string
}

var shortLinks = Links{
	links: make(map[string]string),
}

func main() {
	config.ParseFlags()
	router := gin.Default()
	router.GET("/:id", getURL)
	router.POST("/", postURL)
	err := router.Run(config.FlagPort)
	if err != nil {
		log.Fatal(err)
	}
}

func getURL(c *gin.Context) {
	id := c.Param("id")

	shortLinks.mu.Lock()
	c.Writer.Header().Set("Location", shortLinks.links[id])
	shortLinks.mu.Unlock()

	c.Status(http.StatusTemporaryRedirect)
}

func postURL(c *gin.Context) {
	id := GenerateUniqeString(8)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(400)
		return
	}

	shortLinks.mu.Lock()
	shortLinks.links[id] = string(body)
	shortLinks.mu.Unlock()

	c.Status(http.StatusCreated)
	c.Writer.Header().Set("Content-Type", "text/plain")

	if strings.Contains(config.FlagBaseURL, "http") {
		c.Writer.Write([]byte(config.FlagBaseURL + "/" + id))
		return
	}

	c.Writer.Write([]byte("http://" + config.FlagBaseURL + "/" + id))

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
