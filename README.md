# Backup backed PostgresSQL Image - POC

This images backups and restores the psql image everytime it's shutdown and startup.

## Use cases
- Assume that you have a kubernetes cluster with block devices (like NFS) and the postgres performance on these devices is not as you would like.
- You want to have a postgres image that automatically backups itself

## How does it work?

* Assumption: The postgres database is empty upon startup
* If there is a previous backup file, restore the database
* Start the database
* Upon shutdown, backup the database

* Additionally: create a backup archive every midnight

## Usage

* Assumptions:
    * that a block device is mounted on "/angel" to store the backuped archives
    * that there is enough tempspace "/tmp" and time to backup the database
    

# TODO

* Use "plugins" to support multiple backends (not only postgres)
    * Support simple directories
    * ...?
* Add continuous WAL backup for postgres (point in time recovery)
* How can this be tested?
* Implement scheduled backup

