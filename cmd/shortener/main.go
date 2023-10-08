package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var links = make(map[string]string)

func main() {
	router := gin.Default()
	router.GET("/get/:id", getURL)
	router.POST("/", postURL)
	err := router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}

func postURL(c *gin.Context) {
	id := GenerateUniqeString()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	links[id] = string(body)
	c.String(http.StatusCreated, "http://localhost:8080/"+id+" "+links[id])
}

func getURL(c *gin.Context) {
	c.String(http.StatusTemporaryRedirect, "Location:"+links[c.Param("id")])
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
