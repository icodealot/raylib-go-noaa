package api

import (
	"noaawc/config"

	noaa "github.com/icodealot/noaa"
)

type ForecastPeriod = noaa.ForecastResponsePeriod

type API struct {
	cfg *config.Config
}

func NewAPI(c *config.Config) *API {
	return &API{
		cfg: c,
	}
}

// GetOffice() returns the name of the office responsible for
// the weather observations.
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

// GetForecastPeriods() returns the forecast weather observations
// for the next 14 periods. (Up to 2 per day and up to 8 days)
// When pulling forecast observations in the evening, for example,
// the forecast for "today" may include only 1 observation.
func (api *API) GetForecastPeriods() ([]ForecastPeriod, error) {
	if api.cfg.NOAA.Units == "si" {
		noaa.SetUnits(api.cfg.NOAA.Units)
	}
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
