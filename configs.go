package main

import (
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/phassans/banana/shared"
)

var (
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

	shared.InitLogger()
}
