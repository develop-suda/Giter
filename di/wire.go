// Copyright 2018 The Wire Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The build tag makes sure the stub is not built in the final build.
package di

// import (
// 	"giter/controllers"
// 	"giter/repositories"
// 	"giter/services"

// 	"github.com/google/go-github/github"
// 	"github.com/google/wire"
// 	"github.com/hasura/go-graphql-client"
// 	"gorm.io/gorm"
// )

// func InitCommitRouter(RESTClient *github.Client, GraphQL *graphql.Client) controllers.IRequestController {
// 	wire.Build(controllers.NewRequestController, services.NewRequestService, repositories.NewRequestRepository)
// 	return &controllers.RequestController{}
// }

// func InitAuthRouter(db *gorm.DB) controllers.IAuthControler {
// 	wire.Build(controllers.NewAuthController, services.NewAuthService, repositories.NewAuthRepository)
// 	return &controllers.AuthController{}
// }
