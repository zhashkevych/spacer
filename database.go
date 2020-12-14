package main

import (
	"fmt"
	"os/exec"
)

const (
	pgDumpCommand    = "pg_dump"
	pgRestoreCommand = "pg_restore"
)

// Database is an interface describing DBMS client that creates dump files
type Database interface {
	Dump(file *TempFile) error
}

// Postgres used to dump postgres DB using pg_dump
type Postgres struct {
	Host     string
	Port     string
	Username string
	Name     string
}

// Dump creates dump file with provided name using pg_dump
func (p Postgres) Dump(file *TempFile) error {
	options := p.getExportOptions(file.Name())
	cmd := exec.Command(pgDumpCommand, options...)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) getExportOptions(filename string) []string {
	return []string{
		fmt.Sprintf("-d%s", p.Name),
		fmt.Sprintf("-h%s", p.Host),
		fmt.Sprintf("-p%s", p.Port),
		fmt.Sprintf("-U%s", p.Username),
		"-Ft",
		fmt.Sprintf("-f%s", filename),
	}
}
