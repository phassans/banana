package main

import (
	"net"
	"net/http"
	"time"

	"github.com/phassans/banana/db"
	"github.com/phassans/banana/model"
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/model/favourite"
	"github.com/phassans/banana/model/listing"
	"github.com/phassans/banana/model/notification"
	"github.com/phassans/banana/model/user"
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
	userEngine := user.NewUserEngine(roach.Db, logger)
	businessEngine := business.NewBusinessEngine(roach.Db, logger, userEngine)
	listingEngine := listing.NewListingEngine(roach.Db, logger, businessEngine)
	favouriteEngine := favourite.NewFavoriteEngine(roach.Db, logger, businessEngine, listingEngine)
	notificationEngine := notification.NewNotificationEngine(roach.Db, logger, businessEngine)

	engines := model.NewGenericEngine(
		businessEngine,
		userEngine,
		listingEngine,
		favouriteEngine,
		notificationEngine,
	)

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
