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
	wapiProvider = "WeatherAPI"
	wapiApiBase  = "https://api.weatherapi.com/v1/current.json"
)

type wapiFactory struct {
}

func (f *wapiFactory) Build(client *http.Client) (apis []WeatherApi) {
	coordinates := GetCoordinates("WEX_WAPI_COORDS")
	apiKey := GetStringWithDefault("WEX_WAPI_APIKEY", "")

	for _, coord := range coordinates {
		slog.Info("Creating new WeatherAPI API", "coord", coord)
		apis = append(apis, &wapiApi{client: client, key: apiKey, coord: coord})
	}
	return
}

func init() {
	factories = append(factories, &wapiFactory{})
}

type wapiApi struct {
	client *http.Client
	key    string
	coord  Coordinate
}

func (a *wapiApi) GetCurrentConditions() (*CurrentConditions, error) {
	c, err := a.getCurrent()
	if err != nil {
		return nil, err
	}

	// The WeatherAPI results just include "precipitation", so convert to
	// rain and snow based on the current temperature. This is imperfect,
	// but such is life.
	precipRain, precipSnow := 0.0, 0.0
	if c.Current.TempInC < 0 {
		precipSnow = c.Current.Precip
	} else {
		precipRain = c.Current.Precip
	}

	return &CurrentConditions{
		Provider:     wapiProvider,
		LocationName: c.Location.Name,
		Coordinates:  fmt.Sprintf("%v,%v", a.coord.Lat, a.coord.Lon),
		Description:  c.Current.Condition.Text,
		Temp:         c.Current.TempInC,
		FeelsLike:    c.Current.FeelsLike,
		Humidity:     c.Current.Humidity,
		//PressureGnd:   // This appears to be missing from the Realtime API
		PressureSea:   c.Current.PressureSeaLevel,
		Visibility:    c.Current.Visibility * 1000,
		WindSpeed:     c.Current.WindSpeed * (5.0 / 18.0), // Convert kmph to m/s
		WindDirection: c.Current.WindDir,
		WindGust:      c.Current.WindGust * (5.0 / 18.0), // Convert kmph to m/s
		Clouds:        c.Current.Clouds,
		Rain:          precipRain,
		Snow:          precipSnow,
		UvIndex:       c.Current.UvIndex,
		AqIndex:       c.Current.AirQuality.AqIndex,
		CO:            c.Current.AirQuality.CO,
		NO2:           c.Current.AirQuality.NO2,
		O3:            c.Current.AirQuality.O3,
		SO2:           c.Current.AirQuality.SO2,
		Pm2p5:         c.Current.AirQuality.Pm2p5,
		Pm10:          c.Current.AirQuality.Pm10,
		// NO: // Not Available from WeatherAPI
		// NH3: // Not Available from WeatherAPI
	}, nil
}

type wapiCurrent struct {
	Location struct {
		Name string  `json:"name"`
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
	} `json:"location"`
	Current struct {
		TempInC          float64 `json:"temp_c"`
		FeelsLike        float64 `json:"feelslike_c"`
		Humidity         float64 `json:"humidity"`
		WindSpeed        float64 `json:"wind_kph"`
		WindDir          float64 `json:"wind_degree"`
		WindGust         float64 `json:"gust_kph"`
		Visibility       float64 `json:"vis_km"`
		PressureSeaLevel float64 `json:"presure_mb"`
		Precip           float64 `json:"precip_mm"`
		Clouds           float64 `json:"cloud"`
		UvIndex          float64 `json:"uv"`
		Condition        struct {
			Text string `json:"text"`
			Code int    `json:"code"`
		} `json:"condition"`
		AirQuality struct {
			CO      float64 `json:"co"`
			NO2     float64 `json:"no2"`
			O3      float64 `json:"o3"`
			SO2     float64 `json:"so2"`
			Pm2p5   float64 `json:"pm2_5"`
			Pm10    float64 `json:"pm10"`
			AqIndex float64 `json:"us-epa-index"`
		} `json:"air_quality"`
	} `json:"current"`
}

func (a *wapiApi) getCurrent() (*wapiCurrent, error) {
	url := fmt.Sprintf("%s?key=%s&q=%s&aqi=yes",
		wapiApiBase,
		a.key,
		url.QueryEscape(fmt.Sprintf("%v,%v", a.coord.Lat, a.coord.Lon)),
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
	ret := &wapiCurrent{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
