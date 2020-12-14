package main

import (
	"context"
	"fmt"
)

// Export creates dump and saves it using provided Dumper and Storage objects
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
