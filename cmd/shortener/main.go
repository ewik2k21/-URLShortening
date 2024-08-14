package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-contrib/gzip"

	"github.com/ewik2k21/-URLShortening/cmd/config"
	"github.com/ewik2k21/-URLShortening/internal/app/compressor"
	"github.com/ewik2k21/-URLShortening/internal/app/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Links struct {
	mu    sync.Mutex
	links map[string]string
}

type LinkInput struct {
	URL string `json:"url"`
}

type FileData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var fileData FileData

type LinkOutput struct {
	Result string `json:"result"`
}

var shortLinks = Links{
	links: make(map[string]string),
}

var conn *pgx.Conn

func main() {
	var err error
	if err = initializeLogger(); err != nil {
		panic(err)
	}
	config.ParseFlags()
	conn, err = pgx.Connect(context.Background(), config.FlagConnectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	router := gin.New()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(compressor.DecompressBody())
	router.Use(logger.RequestLogger())
	router.Use(logger.ResponseLogger())
	router.GET("/:id", getURL)
	router.POST("/", postURL)
	router.POST("/api/shorten", postShortenURL)
	router.GET("/ping", getPing)
	err = router.Run(config.FlagPort)
	if err != nil {
		log.Fatal(err)
	}
}

func getPing(c *gin.Context) {

	if err := conn.Ping(context.Background()); err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
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

	AddToFileData(id, linkInput.URL)

	err = WriteDataToFileAsJSON(fileData, config.FlagFileName)
	if err != nil {
		return
	}
}

func getURL(c *gin.Context) {
	id := c.Param("id")

	shortLinks.mu.Lock()
	c.Header("Location", shortLinks.links[id])
	c.Status(http.StatusTemporaryRedirect)
	shortLinks.mu.Unlock()
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
	AddToFileData(id, string(body))

	err = WriteDataToFileAsJSON(fileData, config.FlagFileName)
	if err != nil {
		return
	}
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

func WriteDataToFileAsJSON(input FileData, filedir string) error {

	data, err := json.MarshalIndent(input, "", " ")
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filedir, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return os.WriteFile(filedir, data, 0666)
	}
	file.Write(data)
	return nil
}

func AddToFileData(id string, originalURL string) {
	fileData.ShortURL = id
	fileData.OriginalURL = originalURL
	fileData.UUID = uuid.New().String()
}
