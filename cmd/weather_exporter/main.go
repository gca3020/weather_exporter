package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gca3020/weather_exporter/internal/openweather"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Release string
	SHA     string

	Endpoint = "/metrics"
)

const DefaultPortNumber = 6465

func main() {
	// Create the default logger using logfmt
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	slog.Info("Starting Application", "name", filepath.Base(os.Args[0]), "release", Release, "git", SHA)

	addr, ok := os.LookupEnv("WEX_ADDR")
	if !ok {
		addr = fmt.Sprintf(":%d", DefaultPortNumber)
	}

	prometheus.MustRegister(openweather.NewCollector())
	http.Handle(Endpoint, promhttp.Handler())
	slog.Info("started serving", "addr", addr, "endpoint", Endpoint)

	// Begin serving, which will block forever
	err := http.ListenAndServe(addr, nil)
	slog.Error("ListenAndServe terminated", "err", err)
	os.Exit(1)
}
