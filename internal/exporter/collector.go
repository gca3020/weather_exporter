package exporter

import (
	"log/slog"

	"github.com/gca3020/weather_exporter/internal/api"
	"github.com/prometheus/client_golang/prometheus"
)

var Namespace = "weather"

type Collector struct {
	apis []api.WeatherApi

	description *prometheus.Desc
	temperature *prometheus.Desc
	feelsLike   *prometheus.Desc
	humidity    *prometheus.Desc
	pressureGnd *prometheus.Desc
	pressureSea *prometheus.Desc
	visibility  *prometheus.Desc
	windSpeed   *prometheus.Desc
	windDir     *prometheus.Desc
	windGust    *prometheus.Desc
	clouds      *prometheus.Desc
	rain        *prometheus.Desc
	snow        *prometheus.Desc
	uvi         *prometheus.Desc
	aqi         *prometheus.Desc
	co          *prometheus.Desc
	no          *prometheus.Desc
	no2         *prometheus.Desc
	o3          *prometheus.Desc
	so2         *prometheus.Desc
	nh3         *prometheus.Desc
	pm2p5       *prometheus.Desc
	pm10        *prometheus.Desc
}

func NewCollector(apis []api.WeatherApi) *Collector {
	return &Collector{
		apis: apis,

		description: prometheus.NewDesc(fqName("description"), "Human-readable description of the current conditions", []string{"provider", "location", "desc"}, nil),
		temperature: prometheus.NewDesc(fqName("temperature"), "The temperature at ground level, in Celsius", []string{"provider", "location"}, nil),
		feelsLike:   prometheus.NewDesc(fqName("feelslike"), "The apparent (feels like) temperature at ground level", []string{"provider", "location"}, nil),
		humidity:    prometheus.NewDesc(fqName("humidity"), "The current relative humidity percentage", []string{"provider", "location"}, nil),
		pressureSea: prometheus.NewDesc(fqName("sea_pressure"), "The atmospheric pressure at sea level, in hPa", []string{"provider", "location"}, nil),
		pressureGnd: prometheus.NewDesc(fqName("ground_pressure"), "The atmospheric pressure at the ground level, in hPa", []string{"provider", "location"}, nil),
		visibility:  prometheus.NewDesc(fqName("visibility"), "The visibility, in meters", []string{"provider", "location"}, nil),
		windSpeed:   prometheus.NewDesc(fqName("wind_speed"), "The wind speed, in meters/second", []string{"provider", "location"}, nil),
		windDir:     prometheus.NewDesc(fqName("wind_dir"), "The wind direction, in degrees", []string{"provider", "location"}, nil),
		windGust:    prometheus.NewDesc(fqName("wind_gust"), "The maximum wind gust speed, in meters/second", []string{"provider", "location"}, nil),
		clouds:      prometheus.NewDesc(fqName("cloud_pct"), "The cloud cover percentage", []string{"provider", "location"}, nil),
		rain:        prometheus.NewDesc(fqName("rain"), "The current hourly rainfall rate, in mm", []string{"provider", "location"}, nil),
		snow:        prometheus.NewDesc(fqName("snow"), "The current hourly snowfall rate, in mm", []string{"provider", "location"}, nil),
		uvi:         prometheus.NewDesc(fqName("uv_index"), "The ultraviolet index", []string{"provider", "location"}, nil),
		aqi:         prometheus.NewDesc(fqName("aq_index"), "The air quality index", []string{"provider", "location"}, nil),
		co:          prometheus.NewDesc(fqName("co_conc"), "The carbon monoxide (CO) concentration, in μg/m^3", []string{"provider", "location"}, nil),
		no:          prometheus.NewDesc(fqName("no_conc"), "The nitrogen monoxide (NO) concentration, in μg/m^3", []string{"provider", "location"}, nil),
		no2:         prometheus.NewDesc(fqName("no2_conc"), "The nitrogen dioxide (NO2) concentration, in μg/m^3", []string{"provider", "location"}, nil),
		o3:          prometheus.NewDesc(fqName("o3_conc"), "The ozone (O3) concentration, in μg/m^3", []string{"provider", "location"}, nil),
		so2:         prometheus.NewDesc(fqName("so2_conc"), "The sulfur dioxide (SO2) concentration, in μg/m^3", []string{"provider", "location"}, nil),
		nh3:         prometheus.NewDesc(fqName("nh3_conc"), "The ammonia (NH3) concentration, in μg/m^3", []string{"provider", "location"}, nil),
		pm2p5:       prometheus.NewDesc(fqName("pm2p5_conc"), "The fine particulate (<2.5μm) concentration, in μg/m^3", []string{"provider", "location"}, nil),
		pm10:        prometheus.NewDesc(fqName("pm10_conc"), "The coarse particulate (<10μm) concentration, in μg/m^3", []string{"provider", "location"}, nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	for _, api := range c.apis {
		cc, err := api.GetCurrentConditions()
		slog.Debug("metrics collected", "conditions", cc, "err", err)

		ch <- prometheus.MustNewConstMetric(c.description, prometheus.GaugeValue, 1, cc.Provider, cc.Location, cc.Description)
		ch <- prometheus.MustNewConstMetric(c.temperature, prometheus.GaugeValue, cc.Temp, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.feelsLike, prometheus.GaugeValue, cc.FeelsLike, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.humidity, prometheus.GaugeValue, cc.Humidity, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.pressureSea, prometheus.GaugeValue, cc.PressureSea, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.pressureGnd, prometheus.GaugeValue, cc.PressureGnd, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.visibility, prometheus.GaugeValue, cc.Visibility, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.windSpeed, prometheus.GaugeValue, cc.WindSpeed, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.windDir, prometheus.GaugeValue, cc.WindDirection, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.windGust, prometheus.GaugeValue, cc.WindGust, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.clouds, prometheus.GaugeValue, cc.Clouds, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.rain, prometheus.GaugeValue, cc.Rain, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.snow, prometheus.GaugeValue, cc.Snow, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.uvi, prometheus.GaugeValue, cc.UvIndex, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.aqi, prometheus.GaugeValue, cc.AqIndex, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.co, prometheus.GaugeValue, cc.CO, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.no, prometheus.GaugeValue, cc.NO, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.no2, prometheus.GaugeValue, cc.NO2, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.o3, prometheus.GaugeValue, cc.O3, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.so2, prometheus.GaugeValue, cc.SO2, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.nh3, prometheus.GaugeValue, cc.NH3, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.pm2p5, prometheus.GaugeValue, cc.Pm2p5, cc.Provider, cc.Location)
		ch <- prometheus.MustNewConstMetric(c.pm10, prometheus.GaugeValue, cc.Pm10, cc.Provider, cc.Location)
	}
}

func fqName(name string) string {
	return prometheus.BuildFQName(Namespace, "", name)
}
