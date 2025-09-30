package main

import (
	"flag"
	"fmt"
	"log"

	"go-template/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// migration path, check database/migrations
var migrationPath = "database/migrations"

func main() {
	action := flag.String("action", "up", "migration action: up, down, drop, version")
	steps := flag.Int("steps", 0, "number of steps to migrate (only for up/down)")
	path := flag.String("path", migrationPath, "path to migration files")
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	dbURL := cfg.DB.Config().ConnString()

	// migrate expects: file://<path>
	m, err := migrate.New(
		fmt.Sprintf("file://%s", *path),
		dbURL,
	)
	if err != nil {
		log.Fatalf("failed to init migrate: %v", err)
	}
	defer m.Close()

	switch *action {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	case "drop":
		err = m.Drop()
	case "version":
		version, dirty, verr := m.Version()
		if verr != nil {
			log.Fatalf("failed to get version: %v", verr)
		}
		log.Printf("Current version: %d, Dirty: %v\n", version, dirty)
	default:
		log.Fatalf("unknown action: %s", *action)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("Migration success:", *action)
}
