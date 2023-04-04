package dbmigrate

import (
	"gogql/config"
	"gogql/utils/logger"
	"log"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrateUp(conf config.Config) {
	db, err := sql.Open("postgres", config.MakeDBSource(*conf.DBCreds))
	if err != nil {
		log.Fatal("error when trying to connect:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("db instance is nil:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://sql/migrations", "postgres", driver)
	if err != nil {
		log.Fatal("migration files not found:", err)
	}

	// migrate up
	upErr := m.Up()
	if upErr != nil {
		if upErr == migrate.ErrNoChange {
			log.Println("migration is up to date:", upErr)
		} else {
			log.Fatal("migration up error:", upErr)
		}
	}

	if upErr == nil {
		logger.Success("database migrated up successfully")
	}
}

func RunMigrateDown(conf config.Config) {
	db, err := sql.Open("postgres", config.MakeDBSource(*conf.DBCreds))
	if err != nil {
		log.Fatal("error when trying to connect:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("db instance is nil:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://sql/migrations", "postgres", driver)
	if err != nil {
		log.Fatal("migration files not found:", err)
	}

	// migrate up
	downErr := m.Down()
	if downErr != nil {
		if downErr == migrate.ErrNoChange {
			log.Println("migration is up to date:", downErr)
		} else {
			log.Fatal("migration up error:", downErr)
		}
	}

	if downErr == nil {
		logger.Success("database migrated down successfully")
	}
}
