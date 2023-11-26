package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	rtcache "github.com/ArthurHlt/go-roundtripper-cache"
	"github.com/gca3020/weather_exporter/internal/api"
	"github.com/gca3020/weather_exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Release string
	SHA     string
)

const (
	DefaultAddress = ":9265"
	DefaultTTL     = 10 * time.Minute

	Endpoint = "/metrics"
)

func main() {
	// Create the default logger using logfmt
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	slog.Info("Starting Application", "name", filepath.Base(os.Args[0]), "release", Release, "git", SHA)

	// Get the default parameters that apply to the entire application
	addr := api.GetStringWithDefault("WEX_BIND_ADDR", DefaultAddress)

	// Set up the local client cache with a 10 minute TTL
	ttl := api.GetDurationWithDefault("WEX_TTL", DefaultTTL)
	client := &http.Client{
		Transport: rtcache.NewRoundTripperCache(ttl),
	}

	// Build all the APIs
	apis := api.BuildAll(client)

	// Register the API with the collector and the collector with prometheus.
	prometheus.MustRegister(exporter.NewCollector(apis))
	http.Handle(Endpoint, promhttp.Handler())
	slog.Info("started serving", "addr", addr, "endpoint", Endpoint)

	// Begin serving, which will block forever
	err := http.ListenAndServe(addr, nil)
	slog.Error("ListenAndServe terminated", "err", err)
	os.Exit(1)
}
