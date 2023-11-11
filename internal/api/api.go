package api

type WeatherApi interface {
	GetCurrentConditions() (*CurrentConditions, error)
}

type CurrentConditions struct {
	Provider string // Name of the API Provider (e.g. "OpenWeatherMap", "OpenMeteo", "NOAA")
	Location string // Friendly name of the location to which this conditions apply (e.g. "Denver, US", "Bangkok, Thailand")

	Description   string  // Human-readable description of the current conditions
	Temp          float64 // Temperature at ground level (Celsius)
	FeelsLike     float64 // Apparent, or "feels-like" temperature at ground level (Celsius)
	Humidity      float64 // Relative humidity percent, from 0-100 (percent)
	PressureGnd   float64 // Barometric pressure at ground level (hPa)
	PressureSea   float64 // Barometric pressure at sea level (hPa)
	Visibility    float64 // Visibility (meters)
	WindSpeed     float64 // Wind speed (meters/sec)
	WindDirection float64 // Wind direction (degrees)
	WindGust      float64 // Wind gust speed (meters/sec)
	Clouds        float64 // Cloud cover percentage from 0-100 (percent)
	Rain          float64 // Hourly rainfall rate (mm)
	Snow          float64 // Hourly snowfall rate (mm)
	UvIndex       float64 // The Ultraviolet Index (UVI)
	AqIndex       float64 // The US Air Quality Index (AQI)
	CO            float64 // Carbon Monoxide Concentration (μg/m^3)
	NO            float64 // Nitrogen Monoxide Concentration (μg/m^3)
	NO2           float64 // Nitrogen Dioxide Concentration (μg/m^3)
	O3            float64 // Ozone Concentration (μg/m^3)
	SO2           float64 // Sulfur Dioxide Concentration (μg/m^3)
	NH3           float64 // Ammonia Concentration (μg/m^3)
	Pm2p5         float64 // Fine Particulate Matter (<2.5μm) Concentration (μg/m^3)
	Pm10          float64 // Coarse Particulate Matter (<10μm) Concentration (μg/m^3)
}
