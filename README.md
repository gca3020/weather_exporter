# weather_exporter

A Multi-provider Prometheus Exporter for Weather Metrics

## Motivation

There are a number of existing exporters for Current Weather Conditions, but none that
I could find which supported multiple providers, reported numerous conditions as metrics,
and had containers built for ARM platforms like Raspberry Pi.

## Metrics Reported

| Metric | Description | Notes |
|--------|-------------|-------|
| `weather_description` | Human-readable description of the current conditions | |
| `weather_temperature` | The temperature at ground level, in Celsius | |
| `weather_feelslike` | The apparent (feels like) temperature at ground level | |
| `weather_humidity` | The current relative humidity percentage | |
| `weather_pressure_msl` | The mean atmospheric pressure at sea level (MSL), in hPa | |
| `weather_pressure_surface` | The atmospheric pressure at the ground/surface level, in hPa | |
| `weather_cloud_pct` | The cloud cover percentage | |
| `weather_rain` | The current hourly rainfall rate, in mm | |
| `weather_snow` | The current hourly snowfall rate, in mm | |
| `weather_visibility` | The visibility, in meters | |
| `weather_uv_index` | The ultraviolet index | |
| `weather_wind_dir` | The wind direction, in degrees | |
| `weather_wind_gust` | The maximum wind gust speed, in meters/second | |
| `weather_wind_speed` | The wind speed, in meters/second | |
| `weather_aq_index` | The air quality index | |
| `weather_co_conc` | The carbon monoxide (CO) concentration, in μg/m^3 | |
| `weather_nh3_conc` | The ammonia (NH3) concentration, in μg/m^3 | |
| `weather_no2_conc` | The nitrogen dioxide (NO2) concentration, in μg/m^3 | |
| `weather_no_conc` | The nitrogen monoxide (NO) concentration, in μg/m^3 | |
| `weather_o3_conc` | The ozone (O3) concentration, in μg/m^3 | |
| `weather_pm10_conc` | The coarse particulate (<10μm) concentration, in μg/m^3 | |
| `weather_pm2p5_conc` | The fine particulate (<2.5μm) concentration, in μg/m^3 | |
| `weather_so2_conc` | The sulfur dioxide (SO2) concentration, in μg/m^3 | |

## Configuration

Because it was designed to run in a container, configuration of the weather exporter is
performed via environment variables. Some variables are generic and apply to all providers,
while others are specific to individual providers.

| Variable | Provider | Notes | Default |
|----------|----------|-------|---------|
| `WEX_BIND_ADDR` | All | The address and port that the application binds to | `":6465"`
| `WEX_TTL` | All | To prevent querying remote APIs more frequently than necessary, or exceeding rate limits on API keys, a client side HTTP cache can be enabled. This sets the TTL on the cache | `"10m"` |
| `WEX_OMET_COORDS` | OpenMeteo | Lat/lon pairs for locations to query weather from OpenMeteo, in the format of `"lat,lon;lat,lon"` | `""` |
| `WEX_OW_COORDS` | OpenWeatherMap | Lat/lon pairs for locations to query weather from OpenWeatherMap, in the format of `"lat,lon;lat, lon"` | `""` |
| `WEX_OW_APIKEY` | OpenWeatherMap | The OpenWeatherMap API Key | `""` |
| `WEX_TIO_COORDS` | Tomorrow.io | Lat/lon pairs for locations to query weather from Tomorrow.io | `""` |
| `WEX_TIO_APIKEY` | Tomorrow.io | The Tomorrow.io API Key | `""` |
| `WEX_WAPI_COORDS` | WeatherAPI | Lat/lon pairs for locations to query weather from WeatherAPI.com | `""` |
| `WEX_WAPI_APIKEY` | WeatherAPI | The WeatherAPI API Key | `""` |

## Provider Notes

Though an attempt has been made to normalize the information reported from each provider, there
are some discrepancies between them, which have been captured below.

### OpenWeatherMap

When configured using the `WEX_OW_COORDS` environment variable, the `weather_exporter` will use the
[OpenWeatherMap.org API]([Open](https://openweathermap.org/api)) to retrieve current weather conditions
and make them available to Prometheus. This uses the v2.5 API, which (while out-of-date) does not
require a credit card to sign up.

### Tomorrow.io

The [Tomorrow.io API](https://www.tomorrow.io/weather-api/) provides access to current conditions for
an area using a free API key, but a premium subscription is required for access to the Air Quality API,
so this is not currently implemented.

### Open-Meteo

The [Open-Meteo](https://open-meteo.com/en/docs) API provides for up to 10,000 requests per day without
an API Key required, making it a good free and open source solution. By default, this provider pulls from
multiple weather services, including the NOAA GFS, the ICON, and the European ECMWF.

### WeatherAPI

This configuration uses the [WeatherAPI.com](https://www.weatherapi.com/) Realtime API to query current
conditions for an area, with 15-minute update intervals.

## Examples

### Docker Compose

The following example shows a `docker-compose.yml` file configured to use OpenWeatherMap to report weather metrics in New York City

```yaml
version: "3.8"
services:
  weather_exporter:
    image: gca3020/weather_exporter:latest
    container_name: weather_exporter
    restart: unless-stopped
    ports:
    - '9265:9265'
    environment:
    - WEX_BIND_ADDR=":9265"
    - WEX_TTL="10m"
    - WEX_OW_APIKEY=super-secret-api-key
    - WEX_OW_COORDS="40.75,-73.99"
```

### Prometheus Configuration

```yaml
scrape_configs:
- job_name: "weather"
  scrape_interval: "10m"
  scrape_timeout: "1m"
  static_configs:
  - targets:
    - "weather_exporter:9265"
```
