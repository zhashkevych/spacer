package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	pgDumpCommand    = "pg_dump"
	pgRestoreCommand = "pg_restore"
)

type postgres struct {
	host     string
	port     string
	username string
	dbName   string
	sslMode  string
	password string
}

func (p postgres) dump(filename string) error {
	options := p.getExportOptions(filename)
	cmd := exec.Command(pgDumpCommand, options...)
	err := cmd.Run()
	if err != nil {
		rollback(filename)
		return err
	}

	return nil
}

func (p postgres) getExportOptions(filename string) []string {
	options := []string{
		fmt.Sprintf("-d%s", p.dbName),
		fmt.Sprintf("-h%s", p.host),
		fmt.Sprintf("-p%s", p.port),
		fmt.Sprintf("-U%s", p.username),
		"-Ft",
		fmt.Sprintf("-f%s", filename),
	}

	return options
}

func rollback(filename string) {
	os.Remove(filename)
}
