package main

import (
	"fmt"
	"log"
	"math"
	"noaawc/config"

	rl "github.com/gen2brain/raylib-go/raylib"
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

func main() {

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(800, 450, "Weather Forecast (NOAA)")
	rl.SetTargetFPS(60)

	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("failed to read arguments: %v", err)
	}

	periods, err := getForecastPeriods(cfg)
	if err != nil {
		log.Fatalf("error getting the forecast: %v", err)
	}
	fmt.Printf("%v", periods)

	office := "Chicago, IL (LOT)"
	//office, err := getOfficeName(cfg)
	// if err != nil {
	// 	log.Fatalf("error getting the office name: %v", err)
	// }

	// periodIndex := 0
	// period := periods[periodIndex]
	// fmt.Printf("%v\n", period)

	officePos := rl.NewVector2(0, 20)
	officePos.X = float32(int32(rl.GetScreenWidth())/2 - rl.MeasureText(office, 20)/2)
	tempPeriods := make([]ForecastPeriod, 14)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.DrawText(office, int32(officePos.X), int32(officePos.Y), 20, rl.RayWhite)

		// if rl.IsKeyPressed(rl.KeySpace) {
		// 	periodIndex = (periodIndex + 1) % len(forecast.Periods)
		// 	period = forecast.Periods[periodIndex]
		// 	fmt.Printf("%v\n", period)
		// }

		for i, _ := range tempPeriods {
			// half a sine wave to ease scale in and out
			scale := float32(math.Sin((float64(i) - 0.5) / 14.0 * 180.0))
			if scale < 0.6 {
				continue
			}
			//fmt.Printf("%2.2f\n", scale)
			width := 100 * scale
			xoffset := (100-width)*0.5 - 335
			height := 150 * scale
			yoffset := (150 - height) * 0.5
			padding := 100 / 8

			r := rl.Rectangle{
				X:      float32(padding+i*(100+padding)) + xoffset,
				Y:      280 + yoffset,
				Width:  100 * scale,
				Height: 150 * scale,
			}
			rl.DrawRectangleRoundedLines(r, 0.25, 15, 2*scale, rl.ColorAlpha(rl.RayWhite, scale))

			//rl.DrawRectangleLines(int32(rl.GetScreenWidth()/2), 0, 2.0, int32(rl.GetScreenHeight()), rl.Green)
		}

		// r := rl.Rectangle{X: 10, Y: 280, Width: 100, Height: 150}
		// rl.DrawRectangleRoundedLines(r, 0.25, 15, 2, rl.ColorAlpha(rl.Black, 0.5))

		// rl.DrawText(period.Name, 50, 150, 20, rl.LightGray)
		// rl.DrawText(period.Summary, 50, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

// -----------------------------------------------------------------
// Parse command line arguments and read config file
// if one is provided. Command line arguments include:
//    -config <path to config file>
//    -lat <latitude>
//    -lon <longitude>
// Note: if a config file is provided the lat and lon
// arguments are ignored.
//
// Defaults to Chicago, IL if no arguments are provided.
// -----------------------------------------------------------------
func getConfig() (*config.Config, error) {
	cfg := &config.Config{}
	cfgPath, lat, lon, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("failed to read arguments: %v", err)
	}
	if cfgPath != "" {
		cfg, err = config.NewConfig(cfgPath)
		if err != nil {
			log.Fatalf("failed to parse config file: %v", err)
		}
	} else {
		cfg.NOAA.Latitude = lat
		cfg.NOAA.Longitude = lon
	}
	return cfg, nil
}

func getOfficeName(cfg *config.Config) (string, error) {
	// Grab the cached point so we can identify the Office (CWA)
	p, err := noaa.Points(cfg.NOAA.Latitude, cfg.NOAA.Longitude)
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

func getForecastPeriods(cfg *config.Config) ([]ForecastPeriod, error) {
	// Get the forecast from weather.gov
	forecast, err := noaa.Forecast(cfg.NOAA.Latitude, cfg.NOAA.Longitude)
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
