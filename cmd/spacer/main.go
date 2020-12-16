package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	spacer "github.com/zhashkevych/spacer/pkg"
	"io/ioutil"
	"log"
	"os"
)

const (
	keyPath = "enc.key"
)

func main() {
	parseEnv()
	runApp(initSpacer())
}

func parseEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initSpacer() *spacer.Spacer {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	postgres, err := spacer.NewPostgres(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("failed to create Postgres: %s", err.Error())
	}

	endpoint := os.Getenv("ENDPOINT")
	bucket := os.Getenv("BUCKET")
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	saver, err := spacer.NewSpacesStorage(endpoint, bucket, accessKey, secretKey)
	if err != nil {
		log.Fatalf("failed to create SpacesStorage: %s", err.Error())
	}

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("failed to read key file: %s", err.Error())
	}

	encryptor, err := spacer.NewEncryptor(key)
	if err != nil {
		log.Fatalf("failed to create Encryptor: %s", err.Error())
	}

	var prefix string
	flag.StringVar(&prefix, "prefix", "prod", "dump file prefix (ex. prod/stage)")
	flag.Parse()

	return spacer.NewSpacer(postgres, saver, encryptor, prefix)
}

func runApp(s *spacer.Spacer) {
	app := &cli.App{
		Name: "CLI tool that helps you export encrypted Postgres dumps to DigitalOcean Spaces",
		Commands: []*cli.Command{
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "create and export dump",
				Action:  func(c *cli.Context) error {
					log.Println("Starting export")
					url, err := s.Export(context.Background())
					if err != nil {
						log.Fatal(err)
					}

					log.Println("dump successfully exported to", url)

					return nil
				},
			},
			{
				Name:    "restore",
				Aliases: []string{"r"},
				Usage:   "restore from latest dump",
				Action:  func(c *cli.Context) error {
					log.Println("Starting restore")
					if err := s.Restore(context.Background()); err != nil {
						log.Fatal(err)
					}

					log.Println("DB successfully restored from latest dump")

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}