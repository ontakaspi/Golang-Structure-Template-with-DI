package database

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang-structure-template-with-di/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB Declare the variable for the database
var PostgreDB *gorm.DB

// ConnectDB connect to db
func ConnectDB() {
	var err error

	loggers := logrus.New()
	PgHostAndPort := config.GetEnv("PG_HOST_AND_PORT")
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.GetEnv("PG_USERNAME"),
		config.GetEnv("PG_PASSWORD"),
		PgHostAndPort,
		config.GetEnv("PG_DATABASE_SVC"))

	// Connect to the DB and initialize the DB variable
	PostgreDB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		loggers.Panic("failed to connect database", err.Error())
	}
	loggers.Info("Connection Opened to Database")

}
