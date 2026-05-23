package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/dominicgisler/imap-spam-cleaner/logx"
)

const defaultAddress = ":8080"

func NewHTTPServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/runs/summary", RunSummary)

	return &http.Server{
		Addr:              defaultAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func ListenAndServe() error {
	server := NewHTTPServer()
	defer func() { _ = server.Close() }()
	logx.Infof("Starting HTTP server on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
