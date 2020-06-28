package commands

import (
	"os"
	"os/exec"
	"strconv"
)

type PGRestore struct {

	// pg_restore -h default -U postgres -j 5 -d postgres --data-only -F directory ./blah
	Username   string
	OnlyData   bool
	Format     string
	SourcePath string
	Database   string
	Threads    int
}

func pgRestoreExecutable() string {
	fullCommand, _ := exec.LookPath("pg_restore")
	return fullCommand
}

func (p *PGRestore) Exec() error {
	var args []string

	if len(p.Username) > 0 {
		args = append(args, "--username="+p.Username)
	}
	if p.OnlyData {
		args = append(args, "--data-only")
	}
	if len(p.Format) > 0 {
		args = append(args, "--format="+p.Format)
	}
	if p.Threads > 0 {
		args = append(args, "--jobs="+strconv.Itoa(p.Threads))
	}
	if len(p.Database) > 0 {
		args = append(args, "--dbname="+p.Database)
	}
	if len(p.SourcePath) > 0 {
		args = append(args, p.SourcePath)
	}

	cmd := exec.Command(pgRestoreExecutable(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
