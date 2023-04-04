package server

import (
	"gogql/app/api/handlers"
	"gogql/app/api/routes"
	"gogql/app/master"
	"gogql/app/services"
	"gogql/app/store/dbstore"
	"gogql/app/store/filestore"
	"gogql/config"
)

// All dependency injections will go here
func Injection(c *config.Clients) (*dbstore.DBStore, *routes.Routes) {
	dbs := dbstore.NewDBStore(c.PostgresConn)
	fs := filestore.NewFilestore(c.AWSSession, c.AWSRegion, c.S3BucketName)
	m := master.NewMaster(dbs)
	s := services.NewService(dbs, m)
	h := handlers.NewHandlers(s, fs)
	rt := routes.NewRoutes(h)

	return dbs, rt
}
