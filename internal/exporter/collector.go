package exporter

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

var Namespace = "weather"
var Subsystem = "openweather"

type Collector struct {
	api *Api

	conditions   *prometheus.Desc
	temperature  *prometheus.Desc
	feelsLike    *prometheus.Desc
	pressure     *prometheus.Desc
	humidity     *prometheus.Desc
	pressureSea  *prometheus.Desc
	pressureGrnd *prometheus.Desc
	tempMin      *prometheus.Desc
	tempMax      *prometheus.Desc
	visibility   *prometheus.Desc
	windSpeed    *prometheus.Desc
	windDeg      *prometheus.Desc
	windGust     *prometheus.Desc
	clouds       *prometheus.Desc
	rain1h       *prometheus.Desc
	rain3h       *prometheus.Desc
	snow1h       *prometheus.Desc
	snow3h       *prometheus.Desc
	sunrise      *prometheus.Desc
	sunset       *prometheus.Desc
}

func NewCollector(api *Api) *Collector {
	return &Collector{
		api: api,
		conditions: prometheus.NewDesc(fqName("conditions"),
			"Descriptions of the current conditions, both high-level and detailed", []string{"location", "main", "details"}, nil),
		temperature: prometheus.NewDesc(fqName("temperature"),
			"The current temperature, in degrees Celsius", []string{"location"}, nil),
		feelsLike: prometheus.NewDesc(fqName("feelslike"),
			"The feels like temperature, accounting for human perception of temperature", []string{"location"}, nil),
		pressure: prometheus.NewDesc(fqName("pressure"),
			"The current atmospheric pressure, in hPa", []string{"location"}, nil),
		humidity: prometheus.NewDesc(fqName("humidity"),
			"The current relative humidity percentage", []string{"location"}, nil),
		pressureSea: prometheus.NewDesc(fqName("pressure_sealevel"),
			"The atmospheric pressure at sea level, in hPa", []string{"location"}, nil),
		pressureGrnd: prometheus.NewDesc(fqName("pressure_grndlevel"),
			"The atmospheric pressure at the ground level, in hPa", []string{"location"}, nil),
		tempMin: prometheus.NewDesc(fqName("temp_min"),
			"The minimum temperature in the greater metro area, in degrees Celsius", []string{"location"}, nil),
		tempMax: prometheus.NewDesc(fqName("temp_max"),
			"The maximum temperature in the greater metro area, in degrees Celsius", []string{"location"}, nil),
		visibility: prometheus.NewDesc(fqName("visibility"),
			"The visibility, in meters, up to 10km", []string{"location"}, nil),
		windSpeed: prometheus.NewDesc(fqName("wind_speed"),
			"The wind speed, in meters/second", []string{"location"}, nil),
		windDeg: prometheus.NewDesc(fqName("wind_dir"),
			"The wind direction, in degrees", []string{"location"}, nil),
		windGust: prometheus.NewDesc(fqName("wind_gust"),
			"The maximum wind gust speed, in meters/second", []string{"location"}, nil),
		clouds: prometheus.NewDesc(fqName("cloud_pct"),
			"The cloudiness percent", []string{"location"}, nil),
		rain1h: prometheus.NewDesc(fqName("rain_1h"),
			"Rain volume for the last hour, in millimeters", []string{"location"}, nil),
		rain3h: prometheus.NewDesc(fqName("rain_3h"),
			"Rain volume for the last three hours, in millimeters", []string{"location"}, nil),
		snow1h: prometheus.NewDesc(fqName("snow_1h"),
			"Snow volume for the last hour, in millimeters", []string{"location"}, nil),
		snow3h: prometheus.NewDesc(fqName("snow_3h"),
			"Snow volume for the last three hours, in millimeters", []string{"location"}, nil),
		sunrise: prometheus.NewDesc(fqName("sunrise"),
			"The sunrise time, unix UTC", []string{"location"}, nil),
		sunset: prometheus.NewDesc(fqName("sunset"),
			"The sunset time, unix UTC", []string{"location"}, nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.conditions
	ch <- c.temperature
	ch <- c.feelsLike
	ch <- c.pressure
	ch <- c.humidity
	ch <- c.pressureSea
	ch <- c.pressureGrnd
	ch <- c.tempMin
	ch <- c.tempMax
	ch <- c.visibility
	ch <- c.windSpeed
	ch <- c.windDeg
	ch <- c.windGust
	ch <- c.clouds
	ch <- c.rain1h
	ch <- c.rain3h
	ch <- c.snow1h
	ch <- c.snow3h
	ch <- c.sunrise
	ch <- c.sunset
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	w, err := c.api.getCurrentConditions()
	slog.Debug("metrics collected", "weather", w, "err", err)

	ap, err := c.api.getAirPollution()
	slog.Debug("air pollution", "ap", ap, "err", err)

	uvi, err := c.api.getUvIndex()
	slog.Debug("uv index", "uvi", uvi, "err", err)

	ch <- prometheus.MustNewConstMetric(c.conditions, prometheus.GaugeValue, 1, w.Name, w.Weather[0].Main, w.Weather[0].Description)
	ch <- prometheus.MustNewConstMetric(c.temperature, prometheus.GaugeValue, w.Main.Temp, w.Name)
	ch <- prometheus.MustNewConstMetric(c.feelsLike, prometheus.GaugeValue, w.Main.FeelsLike, w.Name)
	ch <- prometheus.MustNewConstMetric(c.pressure, prometheus.GaugeValue, w.Main.Pressure, w.Name)
	ch <- prometheus.MustNewConstMetric(c.humidity, prometheus.GaugeValue, w.Main.Humidity, w.Name)
	ch <- prometheus.MustNewConstMetric(c.pressureSea, prometheus.GaugeValue, w.Main.SeaLevel, w.Name)
	ch <- prometheus.MustNewConstMetric(c.pressureGrnd, prometheus.GaugeValue, w.Main.GrndLevel, w.Name)
	ch <- prometheus.MustNewConstMetric(c.tempMin, prometheus.GaugeValue, w.Main.TempMin, w.Name)
	ch <- prometheus.MustNewConstMetric(c.tempMax, prometheus.GaugeValue, w.Main.TempMax, w.Name)
	ch <- prometheus.MustNewConstMetric(c.visibility, prometheus.GaugeValue, w.Visibility, w.Name)
	ch <- prometheus.MustNewConstMetric(c.windSpeed, prometheus.GaugeValue, w.Wind.Speed, w.Name)
	ch <- prometheus.MustNewConstMetric(c.windDeg, prometheus.GaugeValue, w.Wind.Deg, w.Name)
	ch <- prometheus.MustNewConstMetric(c.windGust, prometheus.GaugeValue, w.Wind.Gust, w.Name)
	ch <- prometheus.MustNewConstMetric(c.clouds, prometheus.GaugeValue, w.Clouds.All, w.Name)
	ch <- prometheus.MustNewConstMetric(c.rain1h, prometheus.GaugeValue, w.Rain.OneHour, w.Name)
	ch <- prometheus.MustNewConstMetric(c.rain3h, prometheus.GaugeValue, w.Rain.ThreeHour, w.Name)
	ch <- prometheus.MustNewConstMetric(c.snow1h, prometheus.GaugeValue, w.Snow.OneHour, w.Name)
	ch <- prometheus.MustNewConstMetric(c.snow3h, prometheus.GaugeValue, w.Snow.ThreeHour, w.Name)
	ch <- prometheus.MustNewConstMetric(c.sunrise, prometheus.GaugeValue, float64(w.Sys.Sunrise), w.Name)
	ch <- prometheus.MustNewConstMetric(c.sunset, prometheus.GaugeValue, float64(w.Sys.Sunset), w.Name)
}

func fqName(name string) string {
	return prometheus.BuildFQName(Namespace, Subsystem, name)
}
