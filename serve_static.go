package main

import "net/http"

type limitHandler struct {
	connc   chan int
	handler http.Handler
}

func (h *limitHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	select {
	case <-h.connc:
		h.handler.ServeHTTP(w, req)
		h.connc <- 0
	default:
		http.Error(w, "503 too busy", http.StatusServiceUnavailable)
	}
}

func NewLimitHandler(maxConns int, handler http.Handler) http.Handler {
	h := &limitHandler{
		connc:   make(chan int, maxConns),
		handler: handler,
	}
	for i := 0; i < maxConns; i++ {
		h.connc <- i
	}
	return h
}

func InitHandler(handerPrefix, path string, maxConnections int) http.Handler {
	handler := http.FileServer(http.Dir(path))
	limitHandler := NewLimitHandler(maxConnections, handler)

	http.Handle(handerPrefix, http.StripPrefix(handerPrefix, limitHandler))

	return limitHandler
}
