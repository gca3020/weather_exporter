package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	tioProvider = "Tomorrow.io"
	tioApiBase  = "https://api.tomorrow.io/v4/weather/realtime"
)

type tioFactory struct {
}

func (f *tioFactory) Build(client *http.Client) (apis []WeatherApi) {
	coordinates := GetCoordinates("WEX_TIO_COORDS")
	apiKey := GetStringWithDefault("WEX_TIO_APIKEY", "")

	for _, coord := range coordinates {
		slog.Info("Creating new Tomorrow.io API", "coord", coord)
		apis = append(apis, &tioApi{client: client, key: apiKey, coord: coord, units: "metric"})
	}
	return
}

func init() {
	factories = append(factories, &tioFactory{})
}

type tioApi struct {
	client *http.Client
	key    string
	coord  Coordinate
	units  string
}

func (a *tioApi) GetCurrentConditions() (*CurrentConditions, error) {
	c, err := a.getCore()
	if err != nil {
		return nil, err
	}

	return &CurrentConditions{
		Provider:      tioProvider,
		LocationName:  c.Location.Name,
		Coordinates:   fmt.Sprintf("%v,%v", a.coord.Lat, a.coord.Lon),
		Description:   tioCodeToString(c.Data.Values.WeatherCode),
		Temp:          c.Data.Values.Temperature,
		FeelsLike:     c.Data.Values.FeelsLike,
		Humidity:      c.Data.Values.Humidity,
		PressureGnd:   c.Data.Values.PressureSurface,
		PressureSea:   c.Data.Values.PressureSea, // This appears to be missing from the Realtime API
		Visibility:    c.Data.Values.Visibility * 1000,
		WindSpeed:     c.Data.Values.WindSpeed,
		WindDirection: c.Data.Values.WindDirection,
		WindGust:      c.Data.Values.WindGust,
		Clouds:        c.Data.Values.CloudCover,
		Rain:          c.Data.Values.RainIntensity + c.Data.Values.FreezingRainIntensity,
		Snow:          c.Data.Values.SnowIntensity + c.Data.Values.SleetIntensity,
		UvIndex:       c.Data.Values.UvIndex,

		// Air Quality APIs are a Premium subscription, so this is not currently implemented
		//	AqIndex:
		//	CO:
		//	NO:
		//	NO2:
		//	O3:
		//	SO2:
		//	NH3:
		//	Pm2p5:
		//	Pm10:
	}, nil
}

type tioCore struct {
	Data struct {
		Values struct {
			Temperature           float64 `json:"temperature"`
			FeelsLike             float64 `json:"temperatureApparent"`
			CloudCover            float64 `json:"cloudCover"`
			Humidity              float64 `json:"humidity"`
			PressureSurface       float64 `json:"pressureSurfaceLevel"`
			PressureSea           float64 `json:"pressureSeaLevel"`
			RainIntensity         float64 `json:"rainIntensity"`
			FreezingRainIntensity float64 `json:"freezingRainIntensity"`
			SleetIntensity        float64 `json:"sleetIntensity"`
			SnowIntensity         float64 `json:"snowIntensity"`
			WindSpeed             float64 `json:"windSpeed"`
			WindDirection         float64 `json:"windDirection"`
			WindGust              float64 `json:"windGust"`
			Visibility            float64 `json:"visibility"`
			UvIndex               float64 `json:"uvIndex"`
			WeatherCode           int     `json:"weatherCode"`
		} `json:"values"`
	} `json:"data"`
	Location struct {
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
		Name string  `json:"name"`
	} `json:"location"`
}

func (a *tioApi) getCore() (*tioCore, error) {
	url := fmt.Sprintf("%s?location=%s&apikey=%s&units=%s",
		tioApiBase,
		url.QueryEscape(fmt.Sprintf("%v,%v", a.coord.Lat, a.coord.Lon)),
		a.key,
		a.units,
	)
	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		slog.Error("Invalid status code", "code", rsp.StatusCode, "status", rsp.Status)
		return nil, errors.New("invalid response status code")
	}

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	ret := &tioCore{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

var tioCodeMap = map[int]string{
	0:    "Unknown",
	1000: "Clear, Sunny",
	1100: "Mostly Clear",
	1101: "Partly Cloudy",
	1102: "Mostly Cloudy",
	1001: "Cloudy",
	2000: "Fog",
	2100: "Light Fog",
	4000: "Drizzle",
	4001: "Rain",
	4200: "Light Rain",
	4201: "Heavy Rain",
	5000: "Snow",
	5001: "Flurries",
	5100: "Light Snow",
	5101: "Heavy Snow",
	6000: "Freezing Drizzle",
	6001: "Freezing Rain",
	6200: "Light Freezing Rain",
	6201: "Heavy Freezing Rain",
	7000: "Ice Pellets",
	7101: "Heavy Ice Pellets",
	7102: "Light Ice Pellets",
	8000: "Thunderstorm",
}

func tioCodeToString(code int) string {
	if str, ok := tioCodeMap[code]; ok {
		return str
	}
	return fmt.Sprintf("Unknown (%d)", code)
}
