package server

import (
	"gogql/config"
)

func StartApplication(conf config.Config) {
	c := config.SetupClients(conf)
	defer c.PostgresConn.Close()

	restServer := NewRestServer(c)
	restServer.Start(conf.Server.Address)
}
