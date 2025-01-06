package initializer

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

type Clients struct {
	GClient *graphql.Client
	RClient *github.Client
}

func NewGClient() *graphql.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GRAPHQL_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient("https://api.github.com/graphql", httpClient)
	return client
}

// REST API用client作成
func NewRClient() *github.Client {
	return github.NewClient(nil)
}

func NewClients() *Clients {
	return &Clients{GClient: NewGClient(), RClient: NewRClient()}
}
