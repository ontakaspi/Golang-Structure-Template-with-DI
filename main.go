package main

import (
	cors "github.com/itsjamie/gin-cors"
	"golang-structure-template-with-di/config"
	"golang-structure-template-with-di/config/database"
	"golang-structure-template-with-di/router"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

var route *gin.Engine

func init() {
	err := os.Setenv("UNIT_TEST", "0")
	if err != nil {
		return
	}
	// Initialize logger
	config.InitLogFile()
	// Connect to the Database
	database.ConnectDB()
	database.ConnectMongoDB()
	// Do migration and seeding
	database.Migrate()

	initializeBean()

	//setup main routes
	route = gin.New()
	route.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		ValidateHeaders: false,
	}))
	route.GET("/", func(c *gin.Context) {
		pathActive := config.GetEnv("API_PATH_VERSION")
		c.JSON(http.StatusOK, gin.H{"data": "this is not valid endpoint, use active API endpoint ('" + pathActive + "')"})
	})
	route.GET("/health_check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
	//setup version 1 routes
	version := route.Group("/api/v1")

	// route Initialize
	router.InitRoutes(version)
	router.InitRoutesJWT(version)

	// Handler if no route define
	route.NoRoute(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Page not found"})
	})

}
func main() {

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Error(err)
	}
	time.Local = loc // -> this is setting the global timezone

	loggerMain := logrus.New()
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	loggerMain.SetFormatter(formatter)
	loggerMain.Info("Setting global timezone to Asia/Jakarta")

	// Start the cron job
	cronJob()

	p := config.GetEnv("GOLANG_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		loggerMain.Error("Error parse the port")
	}
	if port == 0 {
		loggerMain.Info("APP Port not set, defaulting to 8092")
		port = 8092
	}

	loggerMain.Info("Starting APP Service on Port: " + strconv.Itoa(int(port)))
	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(port)), limit(route))
	if err != nil {
		loggerMain.Info("Error starting APP Service on Port: " + strconv.Itoa(int(port)) + "\n" + err.Error())
		return
	}
}

func cronJob() {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s := gocron.NewScheduler(loc)
	//downloadTrivyDB() every day at 12:00 AM
	go func() {
		s.Every(1).Minute().Do(cleanupVisitors)
	}()

	s.StartAsync()
}
