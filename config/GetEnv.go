package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// EnvVariable use godot package to load/read the .env file and
// return the value of the key
func GetEnv(key string) string {
	var (
		err        error
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)
	unitTest, _ := strconv.Atoi(os.Getenv("UNIT_TEST"))
	if unitTest == 1 {
		replaced := strings.Replace(basepath, "/config", "/", -1)
		replaced = strings.Replace(replaced, "\\config", "\\", -1)
		err = godotenv.Load(replaced + ".env")
	} else {
		err = godotenv.Load(".env")
	}

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
