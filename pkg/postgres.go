package spacer

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	// env variable name used to set postgres password
	pgPassword = "PGPASSWORD"

	// cli tools that are used to dump and restore postgres dbs
	pgDumpCommand    = "pg_dump"
	pgRestoreCommand = "pg_restore"
)

// Postgres used to dump postgres DB using pg_dump
type Postgres struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

func NewPostgres(host, port, username, password, name string) (*Postgres, error) {
	p := &Postgres{Host: host, Port: port, Username: username, Password: password, Name: name}
	if p.Password == "" {
		return p, nil
	}

	err := p.setPasswordForEnv()
	return p, err
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
		"-w",
		"-Ft",
		fmt.Sprintf("-f%s", filename),
	}
}

// setPasswordForEnv helps to disable interactive mode for password input
func (p Postgres) setPasswordForEnv() error {
	return os.Setenv(pgPassword, p.Password)
}