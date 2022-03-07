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

type UICard struct {
	Position       rl.Vector2
	TargetPosition rl.Vector2
	Size           rl.Vector2
	TargetSize     rl.Vector2
	Index          int
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

	office, err := getOfficeName(cfg)
	if err != nil {
		log.Fatalf("error getting the office name: %v", err)
	}

	officePos := rl.NewVector2(0, 20)
	officePos.X = float32(int32(rl.GetScreenWidth())/2 - rl.MeasureText(office, 20)/2)

	cards := createUICards(periods)

	animRunning := false
	animDuration := float32(0.25)
	animTime := float32(0.0)
	deltaTime := float32(0.0)
	selectedIndex := 7

	for !rl.WindowShouldClose() {
		deltaTime = rl.GetFrameTime()

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.DrawText(office, int32(officePos.X), int32(officePos.Y), 20, rl.RayWhite)

		// if rl.IsKeyPressed(rl.KeySpace) {
		// 	periodIndex = (periodIndex + 1) % len(forecast.Periods)
		// 	period = forecast.Periods[periodIndex]
		// 	fmt.Printf("%v\n", period)
		// }

		if rl.IsKeyPressed(rl.KeySpace) && !animRunning {
			// Set new target position for each card.
			animRunning = true
			animTime = 0.0
			tempCard := cards[len(cards)-1] // cache the final card for later
			for i := range cards {

				if cards[i].Index--; cards[i].Index < 0 {
					cards[i].Index = len(cards) - 1
				}

				if i > 0 {
					cards[i].TargetPosition = cards[i-1].Position
					cards[i].TargetSize = cards[i-1].Size
				}
			}
			cards[0].TargetPosition = tempCard.Position
			cards[0].TargetSize = tempCard.Size
			cards[0].Index = tempCard.Index

			// for i := range cards {
			// 	fmt.Printf("%-3s", fmt.Sprint(cards[i].Index))
			// }
			// fmt.Println()
			selectedIndex = (selectedIndex + 1) % len(cards)
			//fmt.Printf("Card %d, Index %d\n", 1, cards[1].Index)
		}

		// Update animation data for the cards.
		if animRunning {
			if animTime += deltaTime; animTime >= animDuration {
				animRunning = false
				for i := range cards {
					cards[i].Position = cards[i].TargetPosition
					cards[i].Size = cards[i].TargetSize
				}
			} else {
				for i, card := range cards {
					cards[i].Position = rl.Vector2Lerp(card.Position, card.TargetPosition, animTime/animDuration)
					cards[i].Size = rl.Vector2Lerp(card.Size, card.TargetSize, animTime/animDuration)
				}
			}
		}

		for _, card := range cards {

			if card.Index < 2 || card.Index > 9 {
				continue
			}

			// half sine wave to ease scale in and out
			//scale := float32(math.Sin((float64(card.Index) - 0.5) / float64(len(cards)) * 180.0))
			scale := sinScale(card.Index, len(cards))
			if scale < 0.5 {
				scale = 0.5
			}
			r := rl.Rectangle{
				X:      card.Position.X,
				Y:      card.Position.Y,
				Width:  card.Size.X,
				Height: card.Size.Y,
			}

			//fmt.Printf("%d\n", card.Index)

			rl.DrawRectangleRounded(r, 0.25, 15, rl.ColorAlpha(rl.LightGray, scale))
			rl.DrawRectangleRoundedLines(r, 0.25, 15, 3*scale, rl.ColorAlpha(rl.RayWhite, scale))
			rl.DrawText(fmt.Sprintf("%d", card.Index+1), int32(card.Position.X+card.Size.X/2), int32(card.Position.Y+card.Size.Y/2), 14, rl.RayWhite)
			//fmt.Printf("Card %d, Alpha %2.2f\n", i, scale)
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

func createUICards(periods []ForecastPeriod) []UICard {
	// Create a slice of UICards
	cards := make([]UICard, 0, len(periods))
	// Create a UICard for each period
	for i := range periods {
		// Create a new UICard
		scale := sinScale(i, len(periods)) //float32(math.Sin((float64(i) - 0.5) / float64(len(periods)) * 180.0))

		// TODO: make UICard sizes configurable
		width := 100 * scale
		xoffset := (100-width)*0.5 - 335 // 14 cards but only 7 on screen at a time shift left to center
		height := 150 * scale
		yoffset := (150 - height) * 0.5
		padding := 100 / 8 // 7 cards on screen + 1 for the final space

		card := UICard{
			Position: rl.Vector2{
				X: float32(padding+i*(100+padding)) + xoffset,
				Y: 280 + yoffset,
			},
			Size: rl.Vector2{
				X: width,
				Y: height,
			},
			Index: i,
		}
		// Add the UICard to the slice
		cards = append(cards, card)
	}
	return cards
}

func sinScale(i int, n int) float32 {
	return float32(math.Sin((float64(i) - 0.5) / float64(n) * 180.0))
}
