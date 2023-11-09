package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/ewik2k21/-URLShortening/cmd/config"
	"github.com/ewik2k21/-URLShortening/internal/app/logger"

	"github.com/gin-gonic/gin"
)

type Links struct {
	mu    sync.Mutex `json:"-"`
	links map[string]string
}

type LinkInput struct {
	URL string `json:"url"`
}

type LinkOutput struct {
	Result string `json:"result"`
}

var shortLinks = Links{
	links: make(map[string]string),
}

func main() {
	if err := initializeLogger(); err != nil {
		panic(err)
	}
	config.ParseFlags()
	router := gin.New()
	router.Use(logger.RequestLogger())
	router.Use(logger.ResponseLogger())
	router.GET("/:id", getURL)
	router.POST("/", postURL)
	router.POST("/api/shorten", postShortenURL)
	err := router.Run(config.FlagPort)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeLogger() error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}
	return nil
}

func postShortenURL(c *gin.Context) {
	id := GenerateUniqeString(8)
	var linkInput LinkInput
	var linkOutput LinkOutput

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		return
	}

	if err = json.Unmarshal(body, &linkInput); err != nil {
		c.Error(err)
		return
	}

	shortLinks.mu.Lock()
	shortLinks.links[id] = linkInput.URL
	shortLinks.mu.Unlock()

	if strings.Contains(config.FlagBaseURL, "http") {
		linkOutput.Result = config.FlagBaseURL + "/" + id
		serializedLink, err := json.MarshalIndent(linkOutput, "", "   ")
		if err != nil {
			c.Error(err)
		}
		c.Data(http.StatusCreated, "application/json", serializedLink)
		return
	}

	linkOutput.Result = "http://" + config.FlagBaseURL + "/" + id
	serializedLink, err := json.MarshalIndent(linkOutput, "", " ")
	if err != nil {
		c.Error(err)
	}
	c.Data(http.StatusCreated, "application/json", serializedLink)

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
