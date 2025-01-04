package main

import (
	"context"
	"fmt"
	"giter/query"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GRAPHQL_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	var query query.GitHubQuery

	err = client.Query(context.Background(), &query, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"commits": query,
		})
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	fmt.Println("server start")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
