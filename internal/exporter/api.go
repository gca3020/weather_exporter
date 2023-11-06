package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiBase = "https://api.openweathermap.org/data/2.5"

type Api struct {
	client   *http.Client
	key      string
	lat, lon float64
	units    string
}

func NewApi(client *http.Client, key string, lat, lon float64) *Api {
	return &Api{
		client: client,
		key:    key,
		lat:    lat,
		lon:    lon,
		units:  "metric",
	}
}

type CurrentConditions struct {
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

type AirPollution struct {
	List []struct {
		Main struct {
			Aqi int64 `json:"aqi"`
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

type UvIndex struct {
	Value float64 `json:"value"`
}

func (a *Api) getCurrentConditions() (*CurrentConditions, error) {
	url := fmt.Sprintf("%s/weather?lat=%f&lon=%f&appid=%s&units=%s", apiBase, a.lat, a.lon, a.key, a.units)
	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	ret := &CurrentConditions{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *Api) getAirPollution() (*AirPollution, error) {
	url := fmt.Sprintf("%s/air_pollution?lat=%f&lon=%f&appid=%s", apiBase, a.lat, a.lon, a.key)
	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	ret := &AirPollution{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *Api) getUvIndex() (*UvIndex, error) {
	url := fmt.Sprintf("%s/uvi?lat=%f&lon=%f&appid=%s", apiBase, a.lat, a.lon, a.key)
	rsp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	ret := &UvIndex{}
	err = json.Unmarshal(rspData, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
