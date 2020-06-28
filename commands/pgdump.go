package commands

import (
	"os"
	"os/exec"
	"strconv"
)

type PGDump struct {
	Username       string
	OnlySchema     bool
	OnlyData       bool
	Format         string
	OutDestination string
	Database       string
	Threads        int
}

func pgDumpExecutable() string {
	fullCommand, _ := exec.LookPath("pg_dump")
	return fullCommand
}

func (p *PGDump) Exec() error {
	var args []string

	if len(p.Username) > 0 {
		args = append(args, "--username="+p.Username)
	}
	if p.OnlySchema {
		args = append(args, "--schema-only")
	}
	if p.OnlyData {
		args = append(args, "--data-only")
	}
	if len(p.Format) > 0 {
		args = append(args, "--format="+p.Format)
	}
	if len(p.OutDestination) > 0 {
		args = append(args, "--file="+p.OutDestination)
	}
	if p.Threads > 0 {
		args = append(args, "--jobs="+strconv.Itoa(p.Threads))
	}
	if len(p.Database) > 0 {
		args = append(args, p.Database)
	}

	cmd := exec.Command(pgDumpExecutable(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
