package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang-structure-template-with-di/config"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB Declare the variable for the database
var MongoDBClient *mongo.Client

func readPreference(mode string) *readpref.ReadPref {
	// Create a ReadPreference from a string mode
	switch mode {
	case "primary":
		return readpref.Primary()
	case "primaryPreferred":
		return readpref.PrimaryPreferred()
	case "secondary":
		return readpref.Secondary()
	case "secondaryPreferred":
		return readpref.SecondaryPreferred()
	case "nearest":
		return readpref.Nearest()
	default:
		return readpref.Primary()
	}
}

// ConnectMongoDB connect to db
func ConnectMongoDB() {
	var err error
	MongoDbHostAndPort := config.GetEnv("MONGO_DB_HOST_AND_PORT")
	MongoDbUser := config.GetEnv("MONGO_DB_USER")
	MongoDbPassword := config.GetEnv("MONGO_DB_PASSWORD")
	//MongoDbDatabase := config.GetEnv("MONGO_DB_NAME")
	ReadPreferenceMongoDB := config.GetEnv("MONGO_DB_READ_PREFERENCE")

	loggers := logrus.New()

	// Connect to the database
	credential := options.Credential{
		Username: MongoDbUser,
		Password: MongoDbPassword,
	}
	clientOptions := options.Client().
		SetReadPreference(readPreference(ReadPreferenceMongoDB)).
		SetHosts(strings.Split(MongoDbHostAndPort, ",")). // Set the host and port
		SetServerSelectionTimeout(30 * time.Second).
		SetAuth(credential)

	MongoDBClient, err = mongo.NewClient(clientOptions)
	if err != nil {
		loggers.Fatalf("failed to setting database MongoDB: %v", err)
	}

	// Menghubungkan ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = MongoDBClient.Connect(ctx)
	if err != nil {
		loggers.Fatalf("failed to connect database MongoDB: %v", err)
	}

	// Memeriksa koneksi MongoDB
	err = MongoDBClient.Ping(ctx, nil)
	if err != nil {
		loggers.Fatalf("failed to ping database MongoDB: %v", err)
	}
	defer cancel()

	loggers.Info("Connection Opened to Database MongoDB")
}
