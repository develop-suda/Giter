// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"giter/controllers"
	"giter/repositories"
	"giter/services"
	"github.com/google/go-github/github"
	"github.com/hasura/go-graphql-client"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitCommitRouter(RESTClient *github.Client, GraphQL *graphql.Client) controllers.IRequestController {
	iRequestRepository := repositories.NewRequestRepository(RESTClient, GraphQL)
	iRequestService := services.NewRequestService(iRequestRepository)
	iRequestController := controllers.NewRequestController(iRequestService)
	return iRequestController
}

func InitAuthRouter(db *gorm.DB) controllers.IAuthControler {
	iAuthRepository := repositories.NewAuthRepository(db)
	iAuthService := services.NewAuthService(iAuthRepository)
	iAuthControler := controllers.NewAuthController(iAuthService)
	return iAuthControler
}
