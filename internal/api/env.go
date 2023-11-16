package api

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Coordinate struct {
	Lat float64
	Lon float64
}

// Parses multiple Lat/Lon pairs, in the format "ENV=12.0,45.0;37.5,109.4"
func GetCoordinates(coordinateEnv string) []Coordinate {
	coordinates := make([]Coordinate, 0)

	// Grab the full list of coordinates from the environment
	coordStr := GetStringWithDefault(coordinateEnv, "")

	if coordStr == "" {
		return nil
	}

	// First split on semicolons to get the coordinate pairs
	coordPairs := strings.Split(strings.TrimSpace(coordStr), ";")
	for _, pair := range coordPairs {
		// Next split on commas to get the lat/long
		tokens := strings.Split(strings.TrimSpace(pair), ",")
		if len(tokens) != 2 && len(tokens) != 0 {
			slog.Error("Coordinate pair does not contain exactly two tokens", "pair", pair)
			continue
		}
		lat, latErr := strconv.ParseFloat(strings.TrimSpace(tokens[0]), 64)
		lon, lonErr := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 64)
		if latErr != nil || lonErr != nil {
			slog.Error("Error parsing latitude/longitude", "latErr", latErr, "lonErr", lonErr)
			continue
		}
		coordinates = append(coordinates, Coordinate{Lat: lat, Lon: lon})
	}
	return coordinates
}

func GetStringWithDefault(env string, defaultVal string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		return defaultVal
	}
	return str
}

func GetDurationWithDefault(env string, defaultVal time.Duration) time.Duration {
	str, ok := os.LookupEnv(env)
	if !ok {
		return defaultVal
	}
	val, err := time.ParseDuration(str)
	if err != nil {
		return defaultVal
	}
	return val
}

func GetIntWithDefault(env string, defaultVal int) int {
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

func GetBoolWithDefault(env string, defaultVal bool) bool {
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
