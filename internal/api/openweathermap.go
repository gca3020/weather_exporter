package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const (
	owmProvider = "OpenWeatherMap"
	owmApiBase  = "https://api.openweathermap.org/data/2.5"
)

type owmFactory struct {
}

func (f *owmFactory) Build(client *http.Client) (apis []WeatherApi) {
	coordinates := GetCoordinates("WEX_OW_COORDS")
	apiKey := GetStringWithDefault("WEX_OW_APIKEY", "")

	for _, coord := range coordinates {
		slog.Info("Creating new OpenWeather API", "coord", coord)
		apis = append(apis, newOwmApi(client, apiKey, coord))
	}
	return
}

func init() {
	factories = append(factories, &owmFactory{})
}

type owmApi struct {
	client *http.Client
	key    string
	coord  Coordinate
	units  string
}

func newOwmApi(client *http.Client, key string, coordinate Coordinate) *owmApi {
	return &owmApi{
		client: client,
		key:    key,
		coord:  coordinate,
		units:  "metric",
	}
}

func (a *owmApi) GetCurrentConditions() (*CurrentConditions, error) {
	c, err := a.getCurrentConditions()
	if err != nil {
		return nil, err
	}
	uv, err := a.getUvIndex()
	if err != nil {
		return nil, err
	}
	ap, err := a.getAirPollution()
	if err != nil {
		return nil, err
	}

	return &CurrentConditions{
		Provider:      owmProvider,
		LocationName:  c.Name,
		Coordinates:   fmt.Sprintf("%v,%v", a.coord.Lat, a.coord.Lon),
		Description:   c.Weather[0].Description,
		Temp:          c.Main.Temp,
		FeelsLike:     c.Main.FeelsLike,
		Humidity:      c.Main.Humidity,
		PressureGnd:   c.Main.GrndLevel, // OWM appears to omit this in the US
		PressureSea:   c.Main.Pressure,  // OWM's primary pressure stat is always present, and is sea-level pressure.
		Visibility:    c.Visibility,
		WindSpeed:     c.Wind.Speed,
		WindDirection: c.Wind.Deg,
		WindGust:      c.Wind.Gust,
		Clouds:        c.Clouds.All,
		Rain:          c.Rain.OneHour,
		Snow:          c.Snow.OneHour,
		UvIndex:       uv.Value,
		AqIndex:       ap.List[0].Main.Aqi,
		CO:            ap.List[0].Components.Co,
		NO:            ap.List[0].Components.No,
		NO2:           ap.List[0].Components.No2,
		O3:            ap.List[0].Components.O3,
		SO2:           ap.List[0].Components.So2,
		NH3:           ap.List[0].Components.Nh3,
		Pm2p5:         ap.List[0].Components.Pm2p5,
		Pm10:          ap.List[0].Components.Pm10,
	}, nil
}

type owCurrentConditions struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  float64 `json:"pressure"`
		Humidity  float64 `json:"humidity"`
		SeaLevel  float64 `json:"sea_level"`
		GrndLevel float64 `json:"grnd_level"`
	} `json:"main"`
	Visibility float64 `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All float64 `json:"all"`
	} `json:"clouds"`
	Rain struct {
		OneHour   float64 `json:"1h"`
		ThreeHour float64 `json:"3h"`
	} `json:"rain"`
	Snow struct {
		OneHour   float64 `json:"1h"`
		ThreeHour float64 `json:"3h"`
	} `json:"snow"`
	Sys struct {
		Sunrise int64 `json:"sunrise"`
		Sunset  int64 `json:"sunset"`
	}
	Name string `json:"name"`
}

type owAirPollution struct {
	List []struct {
		Main struct {
			Aqi float64 `json:"aqi"`
		} `json:"main"`
		Components struct {
			Co    float64 `json:"co"`
			No    float64 `json:"no"`
			No2   float64 `json:"no2"`
			O3    float64 `json:"o3"`
			So2   float64 `json:"so2"`
			Pm2p5 float64 `json:"pm2_5"`
			Pm10  float64 `json:"pm10"`
			Nh3   float64 `json:"nh3"`
		} `json:"components"`
	} `json:"list"`
}

type owUvIndex struct {
	Value float64 `json:"value"`
}

func (a *owmApi) getCurrentConditions() (*owCurrentConditions, error) {
	url := fmt.Sprintf("%s/weather?lat=%f&lon=%f&appid=%s&units=%s", owmApiBase, a.coord.Lat, a.coord.Lon, a.key, a.units)
	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	ret := &owCurrentConditions{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *owmApi) getAirPollution() (*owAirPollution, error) {
	url := fmt.Sprintf("%s/air_pollution?lat=%f&lon=%f&appid=%s", owmApiBase, a.coord.Lat, a.coord.Lon, a.key)
	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	ret := &owAirPollution{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *owmApi) getUvIndex() (*owUvIndex, error) {
	url := fmt.Sprintf("%s/uvi?lat=%f&lon=%f&appid=%s", owmApiBase, a.coord.Lat, a.coord.Lon, a.key)
	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	ret := &owUvIndex{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
