package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gca3020/weather_exporter/internal/openweather"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Release string
	SHA     string
)

const (
	DefaultAddress            = ":6465"
	DefaultOpenweatherEnabled = true

	Endpoint = "/metrics"
)

func main() {
	// Create the default logger using logfmt
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	slog.Info("Starting Application", "name", filepath.Base(os.Args[0]), "release", Release, "git", SHA)

	// Get the default parameters that apply to the entire application
	addr := getStringWithDefault("WEX_BIND_ADDR", DefaultAddress)
	coords, err := getCoordinates("WEX_COORDS")
	if err != nil {
		slog.Error("Could not parse coordinates from environment", "err", err)
		os.Exit(1)
	}

	// Build the OpenWeather API
	client := http.DefaultClient
	api := openweather.NewApi(client, getStringWithDefault("WEX_OW_APIKEY", ""), coords.lat, coords.lon)

	// Register the API with the collector and the collector with prometheus.
	prometheus.MustRegister(openweather.NewCollector(api))
	http.Handle(Endpoint, promhttp.Handler())
	slog.Info("started serving", "addr", addr, "endpoint", Endpoint)

	// Begin serving, which will block forever
	err = http.ListenAndServe(addr, nil)
	slog.Error("ListenAndServe terminated", "err", err)
	os.Exit(1)
}

type coordinate struct {
	lat, lon float64
}

func getCoordinates(coordinateEnv string) (*coordinate, error) {
	coordStr := getStringWithDefault(coordinateEnv, "")
	tokens := strings.Split(coordStr, ",")
	if len(tokens) != 2 {
		return nil, errors.New("coordinate does not contain exactly two fields")
	}
	lat, latErr := strconv.ParseFloat(strings.TrimSpace(tokens[0]), 64)
	lon, lonErr := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 64)
	if latErr != nil || lonErr != nil {
		return nil, errors.Join(errors.New("coordinates are not floating point values"), latErr, lonErr)
	}
	return &coordinate{lat: lat, lon: lon}, nil
}

func getStringWithDefault(env string, defaultVal string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		return defaultVal
	}
	return str
}

/*
func getIntWithDefault(env string, defaultVal int) int {
	str, ok := os.LookupEnv(env)
	if !ok {
		return defaultVal
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultVal
	}
	return val
}

func getBoolWithDefault(env string, defaultVal bool) bool {
	str, ok := os.LookupEnv(env)
	if !ok {
		return defaultVal
	}
	val, err := strconv.ParseBool(str)
	if err != nil {
		return defaultVal
	}
	return val
}
*/
