package api

import "net/http"

var factories []ApiFactory

type ApiFactory interface {
	Build(*http.Client) []WeatherApi
}

func BuildAll(client *http.Client) []WeatherApi {
	apis := make([]WeatherApi, 0)
	for _, factory := range factories {
		apis = append(apis, factory.Build(client)...)
	}
	return apis
}
