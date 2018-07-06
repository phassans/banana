package main

import (
	"net"
	"net/http"
	"time"

	"github.com/phassans/banana/db"
	"github.com/phassans/banana/model"
	"github.com/phassans/banana/route"
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
	favouriteEngine := model.NewFavoriteEngine(roach.Db, logger, businessEngine, listingEngine)

	engines := model.NewGenericEngine(businessEngine, userEngine, listingEngine, favouriteEngine)

	handerPrefix := "/static/"
	path := "./images/"
	maxConnections := 1
	InitHandler(handerPrefix, path, maxConnections)
	go func() { serverErrChannel <- http.ListenAndServe(":3000", nil) }()

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
