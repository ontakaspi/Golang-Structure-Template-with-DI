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
