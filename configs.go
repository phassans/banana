package main

import (
	"net/http"
	"os"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/rs/xlog"
)

var (
	logger          xlog.Logger
	server          http.Server
	serverStartTime time.Time
)

// defaults
var (
	enableDebugLogging = false
	hystrixHTTPTimeout = 30 * time.Second
	maxHTTPConcurrency = 3000
	serverPort         = "8080"
	serverErrChannel   = make(chan error)
)

func config() {
	// record server start time
	serverStartTime = time.Now()

	// Configure hystrix.
	hystrix.DefaultTimeout = int(hystrixHTTPTimeout / time.Millisecond)
	hystrix.DefaultMaxConcurrent = maxHTTPConcurrency

	// set logger level and output format based on env
	level := xlog.LevelInfo
	if enableDebugLogging {
		level = xlog.LevelDebug
	}
	logger = xlog.New(xlog.Config{Level: level, Output: xlog.NewJSONOutput(os.Stdout)})
	xlog.SetLogger(logger)
}
