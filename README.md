# Spacer helps you manage Database dumping/restoring with ease
Spacer provides functionality to dump Postgres database, encrypt and export it to S3-compatible object storage.

Also, it can restore Database using latest saved dump file.

Usage:

```go
package main

import (
	"context"
	"github.com/zhashkevych/spacer/pkg"
	"log"
)

func main() {
	// create DB client
	postgres, err := spacer.NewPostgres("localhost", "5432", "postgres", "qwerty", "postgres")
	if err != nil {
		log.Fatalf("failed to create Postgres: %s", err.Error())
	}

	// create object storage client
	spaces, err := spacer.NewSpacesStorage("ams3.digitaloceanspaces.com", "test-bucket", "yourAccessKey", "yourSecretKey")
	if err != nil {
		log.Fatalf("failed to create SpacesStorage: %s", err.Error())
	}

	// used to encrypt / decrypt dump files
	enc, err := spacer.NewEncryptor([]byte("your-key-should-have-32-bytes!!!"))
	if err != nil {
		log.Fatalf("failed to create encryptor %s", err.Error())
	}

	s := spacer.NewSpacer(postgres, spaces, enc, "prod") // insert prefix in your dump file e.g. prod.dump.*.sql or stage.dump.*.sql

	ctx := context.Background()

	url, err := s.Export(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Dump successfully exported to", url)

	err = s.Restore(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Latest dump successfully restored")
}
``` 

### Use it in your command line:
Prerequisites:
- go 1.15
- pg_dump & pg_restore installed

Steps:
1) Run `make keygen` to generate encryption key
2) Set connection info variables in .env file (look at .env.example)
3) Run `make build` to create binaries
4) Run `./.bin/spacer export` to create dump and export it to storage OR `./bin/spacer restore` to restore DB from latest dump in your storage bucket

## Makefile targets
-  `build` - build spacer binaries
-  `build-keygen` - build keygen binaries
-  `keygen` - build & run keygen ro generate encryption key