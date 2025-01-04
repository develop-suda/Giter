package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			// "commits": commitsQuery,
			// "repositories": repositories,
		})
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	fmt.Println("server start")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
