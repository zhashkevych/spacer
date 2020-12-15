package spacer

import (
	"context"
	"fmt"
)

// Database is an interface describing DBMS client that creates dump files
type Database interface {
	Dump(file *TempFile) error
}

// Storage is used to save/retrive dump file from remote object storage
type Storage interface {
	Save(ctx context.Context, file *TempFile) (string, error)
}

// Export creates dump and saves it using provided Database and Storage objects
func Export(d Database, s Storage) error {
	dumpFile, err := NewTempFile()
	if err != nil {
		return err
	}

	defer dumpFile.Remove()

	if err := d.Dump(dumpFile); err != nil {
		return err
	}

	url, err := s.Save(context.Background(), dumpFile)
	if err != nil {
		return err
	}

	fmt.Println("dump saved at:", url)

	return nil
}
