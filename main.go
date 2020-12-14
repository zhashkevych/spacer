package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var postgres Postgres
	var spacesConfig = SpacesConfig{
		AccessKey: os.Getenv("ACCESS_KEY"),
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	flag.StringVar(&postgres.Host, "db_host", "localhost", "host address")
	flag.StringVar(&postgres.Port, "db_port", "5432", "port value")
	flag.StringVar(&postgres.Username, "db_username", "postgres", "username")
	flag.StringVar(&postgres.Name, "db_name", "postgres", "database name")

	flag.StringVar(&spacesConfig.Endpoint, "endpoint", "localhost", "storage endpoint")
	flag.StringVar(&spacesConfig.Bucket, "bucket", "", "storage bucket name")

	flag.Parse()

	fmt.Println(postgres)
	fmt.Println(spacesConfig)

	saver, err := NewSpacesStorage(spacesConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := Export(postgres, saver); err != nil {
		log.Fatal(err)
	}
}
