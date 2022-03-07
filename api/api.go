package api

import (
	"noaawc/config"

	noaa "github.com/icodealot/noaa"
)

type ForecastPeriod struct {
	ID              int32   `json:"number"`
	Name            string  `json:"name"`
	StartTime       string  `json:"startTime"`
	EndTime         string  `json:"endTime"`
	IsDaytime       bool    `json:"isDaytime"`
	Temperature     float64 `json:"temperature"`
	TemperatureUnit string  `json:"temperatureUnit"`
	WindSpeed       string  `json:"windSpeed"`
	WindDirection   string  `json:"windDirection"`
	Summary         string  `json:"shortForecast"`
	Details         string  `json:"detailedForecast"`
}

type API struct {
	cfg *config.Config
}

func NewAPI(c *config.Config) *API {
	return &API{
		cfg: c,
	}
}

func (api *API) GetOffice() (string, error) {
	// Grab the cached point so we can identify the Office (CWA)
	p, err := noaa.Points(api.cfg.NOAA.Latitude, api.cfg.NOAA.Longitude)
	if err != nil {
		return "", err
	}

	// Get the Office (CWA) name and other useful fields as needed
	o, err := noaa.Office(p.CWA)
	if err != nil {
		return "", err
	}
	return o.Name + " (" + p.CWA + ")", nil
}

func (api *API) GetForecastPeriods() ([]ForecastPeriod, error) {
	// Get the forecast from weather.gov
	forecast, err := noaa.Forecast(api.cfg.NOAA.Latitude, api.cfg.NOAA.Longitude)
	if err != nil {
		return nil, err
	}
	// convert periods []struct to []ForecastPeriod
	periods := make([]ForecastPeriod, 0, len(forecast.Periods))
	for _, period := range forecast.Periods {
		periods = append(periods, period)
	}
	return periods, nil
}
