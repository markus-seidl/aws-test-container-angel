package postgres

import (
	"container-angel/commands"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Postgres struct {
	DataDirectory string
	Username      string
	Password      string
	Port          int
	Database      string
}

func NewPostgres() Postgres {
	return Postgres{
		DataDirectory: dataDirectory(),
		Username:      os.Getenv("POSTGRES_USER"),
		Password:      os.Getenv("POSTGRES_PASSWORD"),
		Port:          5432,
		Database:      "postgres",
	}
}

func dataDirectory() string {
	return os.Getenv("PGDATA")
}

func (p *Postgres) DataDirectoryExists() bool {
	_, err := os.Stat(p.DataDirectory + "/pg_version")
	return err == nil
}

func (p *Postgres) Dump(dumpFilePath string) error {
	schemaDump := commands.PGDump{
		Username:       p.Username,
		OnlySchema:     true,
		OnlyData:       false,
		Format:         "plain",
		OutDestination: dumpFilePath + "/schema",
		Database:       p.Database,
	}
	err := schemaDump.Exec()
	if err != nil {
		log.Fatalf("Schema dump failed: %s\n", err)
		return err
	}

	dataDump := schemaDump
	dataDump.OnlyData = true
	dataDump.OnlySchema = false
	dataDump.Format = "directory"
	dataDump.Threads = 5
	dataDump.OutDestination = dumpFilePath + "/data"

	err = dataDump.Exec()
	if err != nil {
		log.Fatalf("Data dump failed: %s\n", err)
	}
	return err
}

func (p *Postgres) IsAvailable() bool {
	command := "pg_isready"
	fullCommand, _ := exec.LookPath(command)

	cmd := exec.Command(fullCommand)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func (p *Postgres) EnsureAvailable(count int, d time.Duration) {
	success := 0
	for success < count {
		if p.IsAvailable() {
			success += 1
		} else {
			success = 0
		}

		time.Sleep(d)
	}
}

func (p *Postgres) StartPostgres(wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("docker-entrypoint.sh", "postgres")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Postgres failed with %s\n", err)
	}
}
