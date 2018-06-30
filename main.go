package main

import (
	"net"
	"net/http"
	"time"

	"github.com/pshassans/banana/db"
	"github.com/pshassans/banana/model"
	"github.com/pshassans/banana/route"
	"github.com/rs/xlog"
)

func main() {
	// set up defaults and configs
	config()

	// set up DB
	roach, err := db.New(db.Config{Host: "localhost", Port: "5432", User: "pshashidhara", Password: "banana123", Database: "banana"})
	if err != nil {
		xlog.Fatalf("could not connect to db. errpr %s", err)
	}
	xlog.Infof("successfully connected to db")

	// createEngines
	userEngine := model.NewUserEngine(roach.Db, logger)
	businessEngine := model.NewBusinessEngine(roach.Db, logger)
	listingEngine := model.NewListingEngine(roach.Db, logger, businessEngine)

	engines := model.NewGenericEngine(businessEngine, userEngine, listingEngine)

	// start the server
	server = http.Server{Addr: net.JoinHostPort("", serverPort), Handler: route.APIServerHandler(engines)}
	go func() { serverErrChannel <- server.ListenAndServe() }()

	// log server start time
	logger.Infof("API server started at %s. time:%s", server.Addr, serverStartTime)

	// wait for any server error
	select {
	case err := <-serverErrChannel:
		logger.Fatalf("%s service stopped due to error %v with uptime %v", err, time.Since(serverStartTime))
		roach.Close()
	}
}
