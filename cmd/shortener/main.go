package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/gzip"

	"github.com/ewik2k21/-URLShortening/cmd/config"
	"github.com/ewik2k21/-URLShortening/internal/app/compressor"
	"github.com/ewik2k21/-URLShortening/internal/app/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
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

var db *sql.DB

func main() {
	var err error
	if err = initializeLogger(); err != nil {
		panic(err)
	}
	config.ParseFlags()

	router, err := createRouter()
	if err != nil {
		panic(err)
	}

	err = router.Run(config.FlagPort)
	if err != nil {
		log.Fatal(err)
	}
}

func createRouter() (*gin.Engine, error) {
	var err error
	router := gin.New()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(compressor.DecompressBody())
	router.Use(logger.RequestLogger())
	router.Use(logger.ResponseLogger())
	//logger and compressor
	if config.FlagConnectionString != "" {
		fmt.Println(config.FlagConnectionString)
		//db connection and router methods for db
		db, err = sql.Open("pgx", config.FlagConnectionString)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		row := db.QueryRowContext(ctx, "SELECT EXISTS ( SELECT * FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'shorturls') AS table_exists;")

		var tableCheck bool
		err = row.Scan(&tableCheck)

		if err != nil {
			panic(err)
		}
		if !tableCheck {
			//create table
			_, err := db.ExecContext(ctx, "CREATE TABLE shortUrls ("+
				"uuid UUID,"+
				"shortUrl TEXT,"+
				"originalUrl TEXT);")
			if err != nil {
				panic(err)
			}
			fmt.Println("CREATE TABLE ")
		}

		router.GET("/ping", getPing)
		router.GET("/:id", getURL)
		router.POST("/", postURL)
		router.POST("/api/shorten", postShortenURL)
	} else {
		router.GET("/:id", getURL)
		router.POST("/", postURL)
		router.POST("/api/shorten", postShortenURL)
	}
	return router, nil
}

func getPing(c *gin.Context) {

	if err := db.Ping(); err != nil {
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
		c.Error(err)
	}

	if config.FlagConnectionString != "" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err = db.ExecContext(ctx, "INSERT INTO shortsurl (shorturl, originalurl) VALUES ($1, $2);", id, linkInput.URL)
		if err != nil {
			c.Error(err)
		}
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
		c.Error(err)
	}

	if config.FlagConnectionString != "" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err = db.ExecContext(ctx, "INSERT INTO shortsurl (shorturl, originalurl) VALUES ($1, $2);", id, string(body))
		if err != nil {
			c.Error(err)
		}
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
