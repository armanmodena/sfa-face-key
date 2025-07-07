package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var PORT string
var SECURITY_CODE string

var DB_HOST string
var DB_USER string
var DB_PASSWORD string
var DB_NAME string
var DB_PORT string

var MONGO_HOST string
var MONGO_PORT string
var MONGO_DB string
var MONGO_USER string
var MONGO_PASSWORD string

var SFTP_HOST string
var SFTP_USERNAME string
var SFTP_PASSWORD string
var SFTP_PORT string
var SFTP_ROOT string

var JakartaLocation *time.Location

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	PORT = GetEnv("PORT", "9000")
	SECURITY_CODE = GetEnv("SECURITY_CODE", "")

	DB_HOST = GetEnv("DB_HOST", "localhost")
	DB_USER = GetEnv("DB_USER", "postgres")
	DB_PASSWORD = GetEnv("DB_PASSWORD", "")
	DB_NAME = GetEnv("DB_NAME", "")
	DB_PORT = GetEnv("DB_PORT", "5432")

	MONGO_HOST = GetEnv("MONGO_HOST", "192.168.3.86")
	MONGO_PORT = GetEnv("MONGO_PORT", "27017")
	MONGO_DB = GetEnv("MONGO_DB", "sfa_mobile")
	MONGO_USER = GetEnv("MONGO_USER", "root")
	MONGO_PASSWORD = GetEnv("MONGO_PASSWORD", "P@ssw0rd")

	SFTP_HOST = GetEnv("SFTP_HOST", "localhost")
	SFTP_USERNAME = GetEnv("SFTP_USERNAME", "admin")
	SFTP_PASSWORD = GetEnv("SFTP_PASSWORD", "admin")
	SFTP_PORT = GetEnv("SFTP_PORT", "22")
	SFTP_ROOT = GetEnv("SFTP_ROOT", "/upload/arkan/")
}

func InitTimeZone() error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}
	JakartaLocation = loc
	return nil
}
