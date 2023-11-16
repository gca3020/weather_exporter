package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const (
	ometProvider = "Open-Meteo"
)

type ometFactory struct {
}

func (f *ometFactory) Build(client *http.Client) (apis []WeatherApi) {
	coordinates := GetCoordinates("WEX_OMET_COORDS")

	for _, coord := range coordinates {
		slog.Info("Creating new Open-Meteo API", "coord", coord)
		apis = append(apis, newOmetApi(client, coord))
	}
	return
}

func init() {
	factories = append(factories, &ometFactory{})
}

type ometApi struct {
	client *http.Client
	coord  Coordinate
}

func newOmetApi(client *http.Client, coordinate Coordinate) *ometApi {
	return &ometApi{
		client: client,
		coord:  coordinate,
	}
}

func (a *ometApi) GetCurrentConditions() (*CurrentConditions, error) {
	f, err := a.getForecast()
	if err != nil {
		return nil, err
	}
	aq, err := a.getAirQuality()
	if err != nil {
		return nil, err
	}

	return &CurrentConditions{
		Provider:      ometProvider,
		LocationName:  "", // TODO: Reverse Geocoding?
		Coordinates:   fmt.Sprintf("%v,%v", a.coord.Lat, a.coord.Lon),
		Description:   codeToString(f.Current.Code),
		Temp:          f.Current.Temperature,
		FeelsLike:     f.Current.FeelsLike,
		Humidity:      f.Current.Humidity,
		PressureGnd:   f.Current.PressureSurface,
		PressureSea:   f.Current.PressureMsl,
		Visibility:    f.Current.Visibility,
		WindSpeed:     f.Current.WindSpeed,
		WindDirection: f.Current.WindDir,
		WindGust:      f.Current.WindGust,
		Clouds:        f.Current.Clouds,
		Rain:          f.Current.Rain,
		Snow:          f.Current.Snow,
		UvIndex:       aq.Current.Uvi,
		AqIndex:       aq.Current.Aqi,
		CO:            aq.Current.Co,
		NO2:           aq.Current.No2,
		O3:            aq.Current.O3,
		SO2:           aq.Current.SO2,
		NH3:           aq.Current.NH3,
		Pm2p5:         aq.Current.Pm2p5,
		Pm10:          aq.Current.Pm10,
		//NO:          // Not available in Open-Meteo
	}, nil
}

type ometForecast struct {
	Elevation float64 `json:"elevation"`
	Current   struct {
		Temperature     float64 `json:"temperature_2m"`
		Humidity        float64 `json:"relative_humidity_2m"`
		FeelsLike       float64 `json:"apparent_temperature"`
		Rain            float64 `json:"rain"`
		Showers         float64 `json:"showers"`
		Snow            float64 `json:"snowfall"`
		Clouds          float64 `json:"cloud_cover"`
		PressureMsl     float64 `json:"pressure_msl"`
		PressureSurface float64 `json:"surface_pressure"`
		Visibility      float64 `json:"visibility"`
		WindSpeed       float64 `json:"wind_speed_10m"`
		WindDir         float64 `json:"wind_direction_10m"`
		WindGust        float64 `json:"wind_gusts_10m"`
		Code            int     `json:"weather_code"`
	} `json:"current"`
}

type ometAirQuality struct {
	Current struct {
		Uvi   float64 `json:"uv_index"`
		Aqi   float64 `json:"us_aqi"`
		Co    float64 `json:"carbon_monoxide"`
		No2   float64 `json:"nitrogen_dioxide"`
		O3    float64 `json:"ozone"`
		SO2   float64 `json:"sulfur_dioxide"`
		NH3   float64 `json:"ammonia"`
		Pm2p5 float64 `json:"pm2_5"`
		Pm10  float64 `json:"pm10"`
	} `json:"current"`
}

func (a *ometApi) getForecast() (*ometForecast, error) {
	url := fmt.Sprintf("%s?latitude=%v&longitude=%v&current=%s",
		"https://api.open-meteo.com/v1/forecast",
		a.coord.Lat, a.coord.Lon,
		strings.Join([]string{
			"apparent_temperature",
			"cloud_cover",
			"precipitation",
			"pressure_msl",
			"rain",
			"relative_humidity_2m",
			"showers",
			"snowfall",
			"surface_pressure",
			"temperature_2m",
			"weather_code",
			"visibility",
			"wind_direction_10m",
			"wind_gusts_10m",
			"wind_speed_10m",
		}, ","),
	)

	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	ret := &ometForecast{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *ometApi) getAirQuality() (*ometAirQuality, error) {
	url := fmt.Sprintf("%s?latitude=%v&longitude=%v&current=%s",
		"https://air-quality-api.open-meteo.com/v1/air-quality",
		a.coord.Lat, a.coord.Lon,
		strings.Join([]string{
			"ammonia",
			"carbon_monoxide",
			"nitrogen_dioxide",
			"ozone",
			"pm10",
			"pm2_5",
			"sulphur_dioxide",
			"us_aqi",
			"uv_index",
		}, ","),
	)

	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	ret := &ometAirQuality{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

var codeMap = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Fog",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	56: "Light freezing drizzle",
	57: "Dense freezing drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	66: "Light freezing rain",
	67: "Heavy freezing rain",
	71: "Slight snowfall",
	73: "Moderate snowfall",
	75: "Heavy snowfall",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

func codeToString(code int) string {
	if desc, ok := codeMap[code]; ok {
		return desc
	}
	return fmt.Sprintf("Unknown (%d)", code)
}
