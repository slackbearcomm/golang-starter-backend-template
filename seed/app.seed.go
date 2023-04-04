package seed

import (
	"context"
	"gogql/app/master"
	"gogql/app/store/dbstore"
	"gogql/config"
	"gogql/seed/orgseed"
	"gogql/utils/logger"
	"log"
)

func SeedData(conf config.Config) {
	c := config.SetupClients(conf)
	defer c.PostgresConn.Close()

	ctx := context.Background()
	d := dbstore.NewDBStore(c.PostgresConn)
	m := master.NewMaster(d)

	// initiate db transactions
	tx, err := d.DBTX.BeginTx(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer d.DBTX.RollbackTx(ctx, tx)

	// insert admin
	_, err = orgseed.InsertAdmin(tx, m)
	if err != nil {
		log.Fatal(err.Error.Error())
	}

	// insert orgs
	_, err = orgseed.InsertOrganizations(ctx, tx, d, m)
	if err != nil {
		log.Fatal(err.Error.Error())
	}

	// commit transactions to db
	if err := d.DBTX.CommitTx(ctx, tx); err != nil {
		log.Fatal(err)
	}

	logger.Success("Database seeded successfully")
}
