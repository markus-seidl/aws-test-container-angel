package commands

import (
	"os"
	"os/exec"
)

type Psql struct {
	// psql -h default -U postgres -d postgres -f schema
	Username   string
	SourceFile string
	Database   string
}

func pgSqlExecutable() string {
	fullCommand, _ := exec.LookPath("psql")
	return fullCommand
}

func (p *Psql) Exec() error {
	var args []string

	if len(p.Username) > 0 {
		args = append(args, "--username="+p.Username)
	}
	if len(p.SourceFile) > 0 {
		args = append(args, "--file="+p.SourceFile)
	}
	if len(p.Database) > 0 {
		args = append(args, "--dbname="+p.Database)
	}

	cmd := exec.Command(pgSqlExecutable(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
