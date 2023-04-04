package main

import (
	"context"
	"errors"
	"flag"
	"gogql/config"
	"gogql/dbmigrate"
	"gogql/seed"
	"gogql/server"
	"gogql/utils/logger"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

const (
	envPrefix string = "gogql"
)

func main() {
	godotenv.Load(".env.deploy")
	godotenv.Load(".env.local")
	godotenv.Load()

	var (
		rootFlagSet = flag.NewFlagSet("gogql", flag.ExitOnError)
		migrateup   = rootFlagSet.Bool("migrateup", false, "run db up migrations")
		migratedown = rootFlagSet.Bool("migratedown", false, "run db down migrations")
		dbseed      = rootFlagSet.Bool("dbseed", false, "seed the database")
		test        = rootFlagSet.Bool("test", false, "test application")
		run         = rootFlagSet.Bool("run", false, "run the server")
	)

	rootCmd := &ffcli.Command{
		Name:      "root",
		Options:   []ff.Option{ff.WithEnvVarPrefix(envPrefix)},
		ShortHelp: "Run root commands.",
		FlagSet:   rootFlagSet,
		Exec: func(_ context.Context, args []string) error {
			if !*migrateup && !*migratedown && !*dbseed && !*test && !*run {
				return errors.New("-dbseed or -run is required but not provided ")
			}

			if *migrateup {
				migrateUpDatabase()
			}
			if *migratedown {
				migrateDownDatabase()
			}
			if *dbseed {
				seedDatabase()
			}
			if *test {
				testApplication()
			}
			if *run {
				startApplication()
			}
			return nil
		},
	}

	if err := rootCmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		logger.Fatal(err.Error())
	}
}

func migrateUpDatabase() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	dbmigrate.RunMigrateUp(*conf)
}

func migrateDownDatabase() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	dbmigrate.RunMigrateDown(*conf)
}

func seedDatabase() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	seed.SeedData(*conf)
}

func testApplication() {
	logger.Fatal("no test function provided")
}

func startApplication() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	server.StartApplication(*conf)
}
