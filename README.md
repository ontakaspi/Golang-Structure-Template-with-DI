This repository is a golang structured folder that is used to generate a new backend/API service application by using the golang gin framework and use the dependencies injection (https://github.com/goioc/di) light-weight Spring-like library for Go.
# Table of contents
1. [Project Requirment](#project-requirment)
2. [Project Structure](#project-structure)
    1. [Root Directory](#root-directory)
    2. [Route Directory](#route-directory)
    3. [Controller Directory](#controller-directory)
    4. [Service Directory](#service-directory)
    5. [Repository Directory](#repository-directory)
    6. [Model Directory](#model-directory)
    7. [Helper Directory](#helper-directory)
    8. [Middleware Directory](#middleware-directory)
    9. [Library Directory](#library-directory)
    10. [Config Directory](#config-directory)
    11. [Test Directory](#test-directory)
3. [Database Migration](#database-migration)
4. [Run Project](#run-project)

# Project Requirment <a name="project-requirment"></a>

1. Required for this local run in gcc https://jmeubank.github.io/tdm-gcc/
2. Install golang ^v18 here https://go.dev/dl/ and follow the <a href="https://go.dev/doc/install">installation instructions</a>
3. Install postgresql for database transaction (skip if you want to use docker)
4. Install docker (if you want to run this locally)


# Project Structure <a name="project-structure"></a>

### Root Directory <a name="root-directory"></a>
``/`` (root dir)
this folder contain ``.env`` file that store environment data from host, ``_public_key.pem`` for authorize data token JWT (if have jwt Authorization) and file ``main.go`` of The main entrance of the API for setup environment settings, systems,port, etc . For inject the dependencies need add to ``initializeBean.go`` file in the root directory.
<details><summary>Example initializeBean.go</summary>

````go
package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/goioc/di"
	"golang-structure-template-with-di/app/controller"
	"golang-structure-template-with-di/app/repository"
	"golang-structure-template-with-di/app/service"
	"golang-structure-template-with-di/config/database"
	"reflect"
)

func initializeBean() {

	/* GLOBAL - inject the all global */
	_, _ = di.RegisterBeanFactory("dbPostgres", di.Prototype, func(ctx context.Context) (interface{}, error) {
		db := database.PostgreDB
		return db, nil
	})
	_, _ = di.RegisterBeanFactory("validator", di.Singleton, func(ctx context.Context) (interface{}, error) {
		validate := validator.New()
		return validate, nil
	})

	/* CONTROLLER - inject the all controller */
	_, _ = di.RegisterBean("ExampleProjectController", reflect.TypeOf((*controller.ExampleProjectController)(nil)))

	/* SERVICE - inject the all service */
	_, _ = di.RegisterBean("AuthService", reflect.TypeOf((*service.AuthService)(nil)))
	_, _ = di.RegisterBean("ExampleProjectService", reflect.TypeOf((*service.ExampleProjectService)(nil)))

	/* REPOSITORY - inject the all repository */
	_, _ = di.RegisterBean("ExampleProjectRepository", reflect.TypeOf((*repository.ExampleProjectRepository)(nil)))

	/* ----------------- */
	_ = di.InitializeContainer()
}

````
</details>
<details><summary>Example .env</summary>

````dotenv
PUBLIC_KEY_PATH="_public_key.pem"
API_PATH_VERSION="/v1"
PG_HOST=${PG_HOST}
PG_PORT=5432
PG_DATABASE_CONTAINER_SVC='app_db'
PG_USERNAME=${PG_USERNAME}
PG_PASSWORD=${PG_PASSWORD}
PORT=${APP_PORT}


````
</details>
<details><summary>Example main.go</summary>

````go
    package main
        import (
        "fmt"
        "net/http"
        "gopkg.in/gin-gonic/gin.v1"
        "articles/services/mysql"
        "articles/routers/v1"
        "articles/core/models"
    )
    
    var router  *gin.Engine;
    
    func init() {
        mysql.CheckDB()
        router = gin.New();
    router.NoRoute(noRouteHandler())
        version1:=router.Group("/v1")
        v1.InitRoutes(version1)
    
    }
    
    func main() {
        fmt.Println("Server Running on Port: ", 9090)
        http.ListenAndServe(":9090",router)
    }
````
</details>

---
### Route Directory <a name="route-directory"></a>
`/routers` This package will store every routes in your REST API.
The reason separate the handler is, to easy us to manage each routers. So we can create comments about the API , that with apidoc will generate this into structured documentation. Then we will call the function in index.go in current package.
<b>Note, before creating a route, you need to inject the dependency of the controller to the route in InitializeBean.go file in the root directory.</b>
Example:<br>
<details><summary>Example code ExampleRoute.go</summary>

```go
package route

import (
   "github.com/gin-gonic/gin"
   "github.com/goioc/di"
   "golang-structure-template-with-di/app/controller"
)

func SetExampleProjectRoutes(router *gin.RouterGroup) {

   funcControllerExampleProject := di.GetInstance("ExampleProjectController").(*controller.ExampleProjectController)
   router.GET("example", funcControllerExampleProject.ExampleController)

}

```
</details>
<details><summary>Example code routers.go</summary>

```go
package router

import (
	"github.com/gin-gonic/gin"
	"golang-example/app/middleware"
	route "golang-example/router/v1"
)

// InitRoutesJWT function route that use JWT midlleware
func InitRoutesJWT(g *gin.RouterGroup) {
	// Initialize Midlleware
	g.Use(middleware.ErrorHandler())
	g.Use(middleware.JSONMiddleware())
	g.Use(middleware.AuthorizeJWT())
	// Initialize route
	route.ExampleRoute(g)

}

// InitRoutes function route for home or some url not using a JWT Auth
func InitRoutes(g *gin.RouterGroup) {
	g.Use(middleware.ErrorHandler())
	g.Use(middleware.JSONMiddleware())
	// Initialize route
	route.SetHomeRoutes(g)
}

```
</details>

-----
### Controller Directory <a name="controller-directory"></a>
``/app/controllers``
this package will store every controllers in your REST API, and will be used in the routers.
the controller will be used to handle the request and response to the client.
<b> To use the service/repository in the controller(auto-wire), you need to define in the InitializeBean.go file in the root directory and insert the service/repository in the controller with "di.inject" tag.</b>
<details><summary>Example Code</summary>

```go
package controller

import (
   "github.com/gin-gonic/gin"
   "github.com/go-playground/validator/v10"
   "golang-structure-template-with-di/app/helper"
   "golang-structure-template-with-di/app/models/request"
   "golang-structure-template-with-di/app/service"
   "golang-structure-template-with-di/libraries/httpResponse"
)

type ExampleProjectController struct {
   authService                  *service.AuthService                  `di.inject:"AuthService"`
   ExampleProjectService *service.ExampleProjectService `di.inject:"ExampleProjectService"`
   validator                    *validator.Validate                   `di.inject:"validator"`
}

func (global *ExampleProjectController) ExampleController(c *gin.Context) {

   /*-- Check user permission with decoded JWT Token --*/
   checkUserPermission := global.authService.UserHasPermissions(c, "602")
   checkUserRoles := global.authService.UserHasRoles(c, "superadmin")
   if !checkUserRoles && !checkUserPermission {
      httpResponse.Forbidden()
      return
   }

   //*-- Binding and validating data from request body --*/
   var requestData request.CreateExampleProject
   errorValidation := helper.BindAndValidate(c, &requestData, global.validator)
   if errorValidation != nil {
      httpResponse.BadRequest(errorValidation)
   }

   dataResp := global.ExampleProjectService.ExampleProject(requestData)
   httpResponse.HttpCreated(c, "Success Create ContainerProject", dataResp)

}
```


</details>


-----
### Service Directory <a name="service-directory"></a>
``/app/service``
this package will store every services in your REST API, and will be used in the controllers.
the service will create a logic for handling the request and response to the client and pass the data to the controller.
<b> To use the service/repository in the controller(auto-wire), you need to define in the InitializeBean.go file in the root directory and insert the service/repository in the controller with "di.inject" tag.</b>
<details><summary>Example Code</summary>

```go
package service

import (
   "golang-structure-template-with-di/app/models/entity"
   "golang-structure-template-with-di/app/models/request"
   "golang-structure-template-with-di/app/repository"
   "strings"
)

type ExampleProjectService struct {
   ExampleProjectRepository *repository.ExampleProjectRepository `di.inject:"ExampleProjectRepository"`
}

// ExampleProject http request
func (s *ExampleProjectService) ExampleProject(requestContainerProject request.CreateExampleProject) *entity.ExampleProject {

   // Create a new entity.ExampleProject
   //lowercase name request
   projectKey := strings.ToLower(requestContainerProject.Name)
   projectKey = strings.ReplaceAll(projectKey, " ", "-")
   ExampleProject := entity.ExampleProject{
      Key:  projectKey,
      Name: requestContainerProject.Name,
   }
   return s.ExampleProjectRepository.Create(&ExampleProject)
}
```
</details>

-----
### Repository Directory <a name="repository-directory"></a>
``/app/repository``
this package will store every repositories in your REST API, and will be used in the services.
the repository will handling data from services and do the CRUD operation to the database.<b> To use db connection, you need to define the logic to store connection in the InitializeBean.go file in the root directory and insert the service/repository in the controller with "di.inject" tag.</b>
<details><summary>Example Code db postgress</summary>

```go
/* GLOBAL - inject the all global */
_, _ = di.RegisterBeanFactory("dbPostgres", di.Prototype, func(ctx context.Context) (interface{}, error) {
db := database.PostgreDB
return db, nil
})
```
</details>
<details><summary>Example Code</summary>

```go
package repository

import (
   "gorm.io/gorm"
   "golang-structure-template-with-di/app/models/entity"
   "golang-structure-template-with-di/libraries/httpResponse"
)

// ExampleProjectRepository as @Service
type ExampleProjectRepository struct {
   //db *gorm.DB
   db *gorm.DB `di.inject:"dbPostgres"`
}

func (r ExampleProjectRepository) Create(ExampleProject *entity.ExampleProject) *entity.ExampleProject {
   if err := r.db.Create(ExampleProject).Error; err != nil {
      httpResponse.InternalServerError(err)
   }
   return ExampleProject
}

```
</details>

-----
### Models Directory <a name="models-directory"></a>
``/app/models/entity`` This package will store all created model struct for using as data transcaction in database or just as transactional data. we use gorm as ORM library for handling data in database.
<details><summary>Example Code</summary>


```go
package entity

import "gorm.io/gorm"

type ExampleData struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255);not null"`
	Age     int    `gorm:"type:int;not null"`
	Address string `gorm:"type:varchar(255);not null"`
}

```


</details>

``/app/models/request`` This package will store all request as struct data for validate data from API body request.
We use package validator as validation library for handling data from request.
<ul>
<li>
Reference for validator library is <a href="https://pkg.go.dev/github.com/go-playground/validator/v10">https://pkg.go.dev/github.com/go-playground/validator/v10</a>
</li>
<li>
the reference for validating form-data type is <a href="https://github.com/thedevsaddam/govalidator">https://github.com/thedevsaddam/govalidator</a>.
</li>
</ul>

<details><summary>Example Code</summary>

```go
package request
import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// ExampleRequest is a struct for example using validator v10
type ExampleRequest struct {
	Name    string `json:"name" validate:"required,min=3"`
	Age     int    `json:"age" validate:"required,gte=0,lte=130"`
	Address string `json:"address" validate:"required,min=3"`
}

type RequestCreateExampleData struct {
	Name    string
	Age     int
	Address string
}
// CreateExampleData is a struct for example using thedevsaddam/govalidator
type CreateExampleData struct {
	Rules                    map[string][]string
	Message                  map[string][]string
	RequestCreateExampleData RequestCreateExampleData
}

func (std *CreateExampleData) BindRequestField(c *gin.Context) {

	std.Rules = make(map[string][]string)
	std.Rules["name"] = []string{"required"}
	std.Rules["age"] = []string{"required"}
	std.Rules["address"] = []string{"required"}
	std.RequestCreateExampleData.Name = c.PostForm("name")
	std.RequestCreateExampleData.Age, _ = strconv.Atoi(c.PostForm("age"))
	std.RequestCreateExampleData.Address = c.PostForm("address")

}
```


</details>

``/app/models/response`` This package will store all response as struct data for giving API body response.


----
### Helper Directory <a name="helper-directory"></a>
``/app/helper`` This pacakge will store every function that will reusuable in any function in controller.
included helper:
<ul>
<li>
`ChiperHelper.go` helper for encrypting or decrypting data.
</li>
<li>
`ErrorHelper.go` is helper for error handling, reference for error handling is <a href="https://go.dev/blog/error-handling-and-go">https://go.dev/blog/error-handling-and-go</a>
</li>
<li>
`JWTHelper.go` is helper for JWT token, reference for JWT token is <a href="github.com/dgrijalva/jwt-go">github.com/dgrijalva/jwt-go</a>
</li>
<li>
`SimplifyError.go` helper for convert error response package validator <a href="https://pkg.go.dev/github.com/go-playground/validator/v10">v10</a> to human readable.
</li>
<li>
`validate.go` helper for validate data from request.
</li>
</ul>


----
### Middleware Directory <a name="middleware-directory"></a>
``/app/middlewares``
This package will store every middeleware that will use in routes.
included helper:
<ul>
<li>
`errorHandler.go` middleware for error handling that use in routes(gin routes).
</li>
<li>
`JWTMiddleware.go` middleware for JWT authorization that use in routes(gin routes).
</li>
<li>
`JSONMiddleware.go` middleware for giving JSON response that use in routes(gin routes).
</li>
</ul>

-----
### Library Directory <a name="library-directory"></a>
`/libraries`  This package will store any library that used in projects. But only for manually created/imported library, that not available when using go get package_name commands. Could be your own hashing algorithm, graph, tree etc.
include library:
<ul>
<li>
`httpResponse.go` library for giving response to client based on status code and message. The library will give response in JSON format using gin context default like:
`c.JSON(202, SuccessResp{
		Status:  "Progress",
		Message: message
	})
`` and will return panic if some status code error occur. the panic will be handled by errorHandler middleware.
</li>
<li>
`looger.go` library for logging data to file. this package require logrus library.
</li>
</ul>

-----
### Config Directory <a name="config-directory"></a>
``/config`` This package will store any configuration and setting to used in project from any used service, could be mongodb,redis,mysql, elasticsearch, etc.
`/config/database`  This package for database configuration. File `migrations` is for database migration that use package gorm.io/gorm. the migration will be run when project start.

-----
### Test Directory <a name="test-directory"></a>
``/test/mockDatabase`` This package will store any mockDatabase repository.

`/test/tools/`  This package for any tools that used for unit testing, included tools:
<ul>
<li>
`tools.go` a tools for setting driver mock, run test any manymore.
</li>
</ul>
for more information about mock driver, you can see unit testing documentation in UnitTesting.md

-----
# Database Migration <a name="database-migration"></a>
Database migration is a process of creating and updating database tables to match the current model definitions.

1. for using database migration, you need to import package gorm.io/gorm.
2. Create entity struct and use gorm to create table in `/models/entity` folder. for example:
    ```go
    package entity
    
    import "gorm.io/gorm"
    
    type ExampleData struct {
        gorm.Model
        Name    string `gorm:"type:varchar(255);not null"`
        Age     int    `gorm:"type:int;not null"`
        Address string `gorm:"type:varchar(255);not null"`
    }
    ```
   *the `gorm.Model` is a for defining table primary key,created_at,updated_at,deleted_at.*
3. Add this line code to `/config/database/migrations` file function `Migrate()`:
    ```go
    package database
    
    import "golang-example/app/models/entity"
    
    func Migrate() {
        db := PostgreDB
        err := db.AutoMigrate(&entity.ExampleData{})
        if err != nil {
            return
        }
    }
    ```
4. Start project it will run migration automatically.
   *<br>AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes. It will change existing column’s type if its size, precision, nullable changed. It WON’T delete unused columns to protect your data. for more detail about gorm, please refer to https://gorm.io/docs/migration.html*

# Run Project <a name="run-project"></a>
1. Create database in postgresql.
2. Edit `.env` file and set database configuration based on your database configuration.
   example:
   ```dotenv
    PG_HOST=localhost
    PG_PORT=5432
    PG_USERNAME=postgres
    PG_PASSWORD=postgres
    PG_DATABASE_CONTAINER_SVC=postgres
   ```
3. in terminal go to project folder and run command:
    ```shell
    go mod download
    ```
4. then run the project:
    ```shell
    go run main.go
    ```
5. after application run, go to your POSTMAN application to test your API.
