package defaults

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/banana/helper"
	"github.com/rs/xlog"
)

var (
	Logger          xlog.Logger
	Server          http.Server
	ServerStartTime time.Time
	Ctx             context.Context
)

// defaults
var (
	CurrentApiVersion  = "/v1"
	EnableDebugLogging = false
	HystrixHTTPTimeout = 3 * time.Second
	MaxHTTPConcurrency = 3000
	ServerPort         = "8080"
	ServerErrChannel   = make(chan error)
)

func Config() {
	// record server start time
	ServerStartTime = time.Now()

	// Configure hystrix.
	hystrix.DefaultTimeout = int(HystrixHTTPTimeout / time.Millisecond)
	hystrix.DefaultMaxConcurrent = MaxHTTPConcurrency

	// set logger level and output format based on env
	level := xlog.LevelInfo
	if EnableDebugLogging {
		level = xlog.LevelDebug
	}
	Logger = xlog.New(xlog.Config{Level: level, Output: xlog.NewJSONOutput(os.Stdout)})
	xlog.SetLogger(Logger)

	Ctx = helper.NewContext()
}
