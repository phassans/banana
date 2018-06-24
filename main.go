package main

import (
	"net"
	"net/http"
	"time"

	"github.com/banana/database"
	"github.com/banana/routes"
)

func main() {
	// set up defaults and configs
	config()

	// set up DB
	database.InitializeDatabase()

	// start the server
	server = http.Server{Addr: net.JoinHostPort("", serverPort), Handler: routes.APIServerHandler()}
	go func() { serverErrChannel <- server.ListenAndServe() }()

	// log server start time
	logger.Infof("API server started at %s. time:%s", server.Addr, serverStartTime)

	// wait for any server error
	select {
	case err := <-serverErrChannel:
		logger.Fatalf("%s service stopped due to error %v with uptime %v", err, time.Since(serverStartTime))
	}
}
