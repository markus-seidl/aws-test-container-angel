package main

import (
	"container-angel/angel"
	"container-angel/commands"
	"container-angel/postgres"
	"container-angel/store"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	println("Container angel started.")
	conf := angel.DefaultConfiguration()

	_ = os.Remove("/rdy") // ensure that no rdy file exists for the live/readiness check

	db := postgres.NewPostgres()

	skipRestore := false
	if db.DataDirectoryExists() {
		println("Data directory already exists, restoration will be skipped!")
		skipRestore = true
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go db.StartPostgres(&wg)

	db.EnsureAvailable(10, time.Second)

	ListenForSignals(conf, db)

	// Restore database if needed
	backupArchives := store.FindAll(conf.AngelDirectory)
	if len(backupArchives.Archives) > 0 && !skipRestore {
		restoreArchive := backupArchives.Archives[0]
		log.Printf("Restoring %s\n", restoreArchive)

		Restore(conf, db, restoreArchive)
	} else {
		log.Println("No archive to restore / skipped restore.")
	}

	_, _ = os.OpenFile("/rdy", os.O_RDONLY|os.O_CREATE, 0666)
	log.Println("Database available.")

	wg.Wait()
	println("Container angel ended.")
}

func Restore(conf angel.Configuration, db postgres.Postgres, archive store.Archive) {
	tempDirectory, err := ioutil.TempDir(conf.TempDirectory, "angel")
	if err != nil {
		log.Fatal(err)
	}
	err = archiver.Unarchive(archive.FQP, tempDirectory)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDirectory)

	// Restore schema
	schemaRestore := commands.Psql{
		Username:   "postgres",
		SourceFile: tempDirectory + "/schema",
		Database:   "postgres",
	}
	err = schemaRestore.Exec()
	if err != nil {
		log.Fatalf("Restore schema produced error %s\n", err)
	}
	log.Println("Schema restored.")

	// Restore data
	dataRestore := commands.PGRestore{
		Username:   "postgres",
		OnlyData:   true,
		Format:     "directory",
		SourcePath: tempDirectory + "/data",
		Database:   "postgres",
		Threads:    5,
	}
	err = dataRestore.Exec()
	if err != nil {
		log.Fatalf("Restore schema produced error %s\n", err)
	}
	log.Println("Data restored.")
}

func Backup(conf angel.Configuration, db postgres.Postgres, archive store.Archive) {
	tempDirectory, err := ioutil.TempDir(conf.TempDirectory, "angel")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDirectory)

	log.Printf("Starting backup %s...\n", archive.FQP)
	err = db.Dump(tempDirectory)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare input for archiver (have to iterate ourself, the archiver has problems with directories)
	files, err := ioutil.ReadDir(tempDirectory)
	if err != nil {
		log.Fatal(err)
	}

	var archivePackage []string
	for _, element := range files {
		archivePackage = append(archivePackage, tempDirectory+"/"+element.Name())
	}
	err = archiver.Archive(archivePackage, archive.FQP+".tar.gz")

	log.Println("...finished.")
}

func ListenForSignals(conf angel.Configuration, db postgres.Postgres) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go HandleSignal(conf, db, signalChannel)
}

func HandleSignal(conf angel.Configuration, db postgres.Postgres, c <-chan os.Signal) {
	for true {
		s := <-c
		log.Println("Got Signal", s)

		archive := store.NewArchive(conf)
		Backup(conf, db, archive)

		os.Exit(0)
	}
}
