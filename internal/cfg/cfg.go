package cfg

import (
	"log"
	"os"
	"strconv"
)

type Cfg struct {
	DBName         string
	DBPort         int
	DBHost         string
	DBUsername     string
	DBPassword     string
	PrometheusAddr string
	JWTSecret      string
	BCryptSalt     int
	S3ID           string
	S3SecretKey    string
	S3BucketName   string
	S3Region   string
}

func Load() *Cfg {
	var err error
	cfg := &Cfg{}

	cfg.DBName = os.Getenv("DB_NAME")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBUsername = os.Getenv("DB_USERNAME")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.DBPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("fail convert db port to int:", err)
	}

	return cfg
}
