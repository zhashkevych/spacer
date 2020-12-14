package main

import (
	"fmt"
	"log"
	"time"
)

/*

Project Description:
- CLI tool for quick Postgres dumps creation and following export to DigitalOcean Spaces (using encryption)

Export cmd:
	pg_dump -d test -h localhost -p 5436 -U postgres -Ft > backup_test.tar.gz

Restore cmd:
	pg_restore -d test -p 5436 -h localhost  backup_test.tar.gz -c -U postgres

*/

const filenameTemplate = "dump_%d.tar.gz"

func main() {
	pg := postgres{
		host:     "localhost",
		port:     "5436",
		dbName:   "postgres",
		username: "postgres",
	}

	filename := generateFilename()
	if err := pg.dump(filename); err != nil {
		log.Fatalf("export failure: %s", err.Error())
	}
}

func generateFilename() string {
	return fmt.Sprintf(filenameTemplate, time.Now().Unix())
}

// func newPostgres(host, port, username, db) postgres {
// 	var cfg pgConfig

// 	flag.StringVar(&cfg.Host, "host", "localhost", "host address")
// 	flag.StringVar(&cfg.Port, "port", "5432", "port value")
// 	flag.StringVar(&cfg.Username, "username", "postgres", "username")
// 	flag.StringVar(&cfg.Password, "password", "qwerty", "password")
// 	flag.StringVar(&cfg.DBName, "dbname", "postgres", "db name")
// 	flag.StringVar(&cfg.SSLMode, "sslmode", "disable", "sslmode")

// 	return cfg
// }
