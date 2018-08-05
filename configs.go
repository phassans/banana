package main

import (
	"net/http"
	"os"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/rs/zerolog"
)

var (
	logger          zerolog.Logger
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

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("role", "my-service").
		Str("host", "hungryhour").
		Logger()
	logger.Info().Msg("please json")

}
