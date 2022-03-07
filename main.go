package main

import (
	"log"
	"noaawc/api"
	"noaawc/config"
	"noaawc/ui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// TODO: add configuration options for the UI
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to read arguments: %v", err)
	}

	API := api.NewAPI(cfg)

	periods, err := API.GetForecastPeriods()
	if err != nil {
		log.Fatalf("error getting the forecast: %v", err)
	}
	//fmt.Printf("%v\n", periods)

	office, err := API.GetOffice()
	if err != nil {
		log.Fatalf("error getting the office name: %v", err)
	}

	cards := ui.NewUIScrollingCards(len(periods), 0.25)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(800, 450, "Weather Forecast (NOAA)")
	rl.SetTargetFPS(60)

	officePos := rl.NewVector2(
		float32(int32(rl.GetScreenWidth())/2-rl.MeasureText(office, 20)/2),
		20,
	)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeySpace) {
			cards.BeginScrolling() // ignored if already scrolling
		}
		cards.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.DrawText(office, int32(officePos.X), int32(officePos.Y), 20, rl.RayWhite)

		// TODO: create interface for card image thumbnails and
		// add those thumbnail images to the scrolling cards

		// if rl.IsKeyPressed(rl.KeySpace) {
		// 	periodIndex = (periodIndex + 1) % len(forecast.Periods)
		// 	period = forecast.Periods[periodIndex]
		// 	fmt.Printf("%v\n", period)
		// }

		cards.Draw()

		// r := rl.Rectangle{X: 10, Y: 280, Width: 100, Height: 150}
		// rl.DrawRectangleRoundedLines(r, 0.25, 15, 2, rl.ColorAlpha(rl.Black, 0.5))

		// rl.DrawText(period.Name, 50, 150, 20, rl.LightGray)
		// rl.DrawText(period.Summary, 50, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
