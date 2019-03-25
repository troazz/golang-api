package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"./stores"
	"./structs"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const limit = 10

var cc *cache.Cache

func main() {
	cc = cache.New(5*time.Minute, 10*time.Minute)
	go http()
	worker()
}

func worker() {
	msgs, close, err := stores.Subscribe("news")
	if err != nil {
		panic(err)
	}
	defer close()

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			stores.Save(d.Body)
			cc.Flush()
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func http() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3001"
	}

	r := gin.Default()

	r.GET("/news", GetNews)
	r.POST("/news", PostNews)

	r.Run()
}

func GetNews(c *gin.Context) {
	var err error
	var errors []string
	var page int
	var allnews []structs.News
	var key string

	// get page request from query string
	q := c.Request.URL.Query()
	page, err = strconv.Atoi(strings.Join(q["page"], ""))
	if page <= 0 {
		page = 1
	}

	key = "news" + string(page) + string(limit)
	if x, found := cc.Get(key); found {
		fmt.Println("From CACHE")
		allnews = x.([]structs.News)
	} else {
		// get data from store
		allnews, err = stores.Get(page, limit)
		if err != nil {
			errors = append(errors, err.Error())
		}
		fmt.Println("SET CACHE")
		cc.Set(key, allnews, cache.DefaultExpiration)
	}

	if len(errors) > 0 {
		c.JSON(400, gin.H{"errors": errors})
	} else {
		c.JSON(200, gin.H{
			"page":  page,
			"limit": limit,
			"data":  allnews,
		})
	}
}

func PostNews(c *gin.Context) {
	var news structs.News
	var errors []string

	author := c.PostForm("author")
	body := c.PostForm("body")

	if author == "" {
		errors = append(errors, "Author is required")
	}
	if body == "" {
		errors = append(errors, "Body is required")
	}

	if len(errors) > 0 {
		c.JSON(400, gin.H{"errors": errors})
	} else {
		news.Author = author
		news.Body = body

		if err := stores.Queue("news", news); err != nil {
			errors = append(errors, err.Error())
			c.JSON(400, gin.H{"errors": errors})
			panic(err)
		}

		c.JSON(200, gin.H{"message": "News queued for saving"})
	}
}
