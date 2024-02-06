package router

import (
	"github.com/gin-gonic/gin"
	"golang-structure-template-with-di/app/middleware"
	route "golang-structure-template-with-di/router/v1"
)

func InitRoutesJWT(g *gin.RouterGroup) {
	// Initialize Midlleware
	g.Use(middleware.ErrorHandler())
	g.Use(middleware.JSONMiddleware())

	g.Use(func(c *gin.Context) {
		// Define a list of endpoints that should be exempt from JWT authorization
		exemptedEndpoints := []string{
			"/api/v1/integration/example/scan", // Add your exempted endpoints here
			"/api/v1/integration/example/:image_id",
		}

		// Get the current request path
		requestPath := c.FullPath()

		// Check if the current endpoint is exempted
		for _, endpoint := range exemptedEndpoints {
			if requestPath == endpoint {
				// Endpoint is exempted, don't perform JWT authorization
				return
			}
		}
		// For all other endpoints, perform JWT authorization
		middleware.AuthorizeJWT()(c)
	})
	// Initialize route
	route.SetExampleProjectRoutes(g)
	//route.SetImageScanRoutes(g)
	//route.SetConfigurationScanRoutes(g)

}

// InitRoutes function route for home or some url not using a JWT Auth
func InitRoutes(g *gin.RouterGroup) {
	g.Use(middleware.ErrorHandler())
	g.Use(middleware.JSONMiddleware())

	// Initialize route
	route.SetHomeRoutes(g)
}
