package config

import (
	"os"
	"strconv"
	"time"

	"fampay-assignment/logger"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	Port                     string
	DataDbHost               string
	DataDbUser               string
	DataDbPassword           string
	RedisUri                 string
	NewRelicKey              string
	DataDbPasswordSecretName string
	AllowedOrigins           []string
	DataDbPort               int
	YoutubeApiKey1           string
	YoutubeApiKey2           string
	YoutubeApiKey3           string
)

var (
	QUERY_TIMEOUT        = 10 * time.Second
	REDIS_TIMEOUT        = 1 * time.Second
	YOUTUBE_SEARCH_QUERY = "news"
	MAX_PAGINATION_SIZE  = 10
	DATE_FORMAT          = "2006-01-02T15:04:05Z"
	CORS_ALLOWED_METHODS = []string{"GET", "POST", "PATCH", "OPTIONS"}
	CORS_ALLOWED_HEADERS = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Accept",
		"Cookie",
	}
	CACHE_TTL = 5 * time.Minute
)

func mustGetEnvVar(name string) string {
	val, exists := os.LookupEnv(name)
	if !exists {
		logger.Log.WithFields(logrus.Fields{
			"name": name,
		}).Fatal("env var not found")
	}
	return val
}

func parseEnvs() {
	err := godotenv.Load()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"fileName": ".env",
		}).Warn("failed to load env file")
	}

	Port = mustGetEnvVar("PORT")
	RedisUri = mustGetEnvVar("REDIS_URI")
	dataDbPort := mustGetEnvVar("DATA_DB_PORT")
	YoutubeApiKey1 = mustGetEnvVar("YOUTUBE_API_KEY1")
	YoutubeApiKey2 = mustGetEnvVar("YOUTUBE_API_KEY2")
	YoutubeApiKey3 = mustGetEnvVar("YOUTUBE_API_KEY3")
	DataDbHost = mustGetEnvVar("DATA_DB_HOST")
	DataDbPassword = mustGetEnvVar("DATA_DB_PASSWORD")
	DataDbUser = mustGetEnvVar("DATA_DB_USER")
	DataDbPort, err = strconv.Atoi(dataDbPort)
	AllowedOrigins = []string{"*"}
	if err != nil {
		logger.Log.WithField("port", dataDbPort).Fatal("invalid data db port")
	}

}

func init() {
	parseEnvs()
}
