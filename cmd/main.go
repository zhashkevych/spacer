package main

import (
	"github.com/joho/godotenv"
	"github.com/zhashkevych/spacer/pkg"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	postgres, err := spacer.NewPostgres(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal(err)
	}

	endpoint := os.Getenv("ENDPOINT")
	bucket := os.Getenv("BUCKET")
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	saver, err := spacer.NewSpacesStorage(endpoint, bucket, accessKey, secretKey)
	if err != nil {
		log.Fatal(err)
	}

	if err := spacer.Export(postgres, saver); err != nil {
		log.Fatal(err)
	}
}
