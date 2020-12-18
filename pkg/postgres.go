package spacer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

const (
	// env variable name used to set postgres password
	pgPassword = "PGPASSWORD"

	// cli tools that are used to dump and restore postgres dbs
	pgDump    = "pg_dump"
	pgRestore = "pg_restore"
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

	// If password provided, we need to set PGPASSWORD env var
	return p, p.setPasswordForEnv()
}

// Dump creates dump file with provided name using pg_dump
func (p Postgres) Dump(ctx context.Context, filename string) error {
	options := p.getDumpOptions(filename)
	cmd := exec.CommandContext(ctx, pgDump, options...)

	return cmd.Run()
}

func (p Postgres) getDumpOptions(filename string) []string {
	return []string{
		fmt.Sprintf("-d%s", p.Name),
		fmt.Sprintf("-h%s", p.Host),
		fmt.Sprintf("-p%s", p.Port),
		fmt.Sprintf("-U%s", p.Username),
		"-Ft",
		fmt.Sprintf("-f%s", filename),
	}
}

func (p Postgres) Restore(ctx context.Context, filename string) error {
	options := p.getRestoreOptions(filename)
	cmd := exec.CommandContext(ctx, pgRestore, options...)

	return cmd.Run()
}

func (p Postgres) getRestoreOptions(filename string) []string {
	return []string{
		fmt.Sprintf("-d%s", p.Name),
		fmt.Sprintf("-h%s", p.Host),
		fmt.Sprintf("-p%s", p.Port),
		fmt.Sprintf("-U%s", p.Username),
		"-Ft",
		filename,
	}
}

// setPasswordForEnv helps to disable interactive mode for password input
func (p Postgres) setPasswordForEnv() error {
	return os.Setenv(pgPassword, p.Password)
}
