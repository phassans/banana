package main

import (
	"net"
	"net/http"
	"time"

	"github.com/banana/database"
	"github.com/banana/defaults"
	"github.com/banana/routes"
)

func main() {
	// set up defaults and configs
	defaults.Config()

	// set up DB
	database.InitializeDatabase()

	// start the server
	defaults.Server = http.Server{Addr: net.JoinHostPort("", defaults.ServerPort), Handler: routes.WebServerRouter(defaults.Ctx)}
	go func() { defaults.ServerErrChannel <- defaults.Server.ListenAndServe() }()

	// log server start time
	defaults.Logger.Infof("API server started at %s. time:%s", defaults.Server.Addr, defaults.ServerStartTime)

	// wait for any server error
	select {
	case err := <-defaults.ServerErrChannel:
		defaults.Logger.Fatalf("%s service stopped due to error %v with uptime %v", err, time.Since(defaults.ServerStartTime))
	}
}
