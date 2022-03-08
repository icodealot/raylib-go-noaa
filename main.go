package main

import (
	"fmt"
	"log"
	"noaawc/api"
	"noaawc/config"
	"noaawc/ui"
	"strings"

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

	// Note: if you need to debug the weather data you uncomment these lines
	// and check the output in the console.
	// ------------------------------------------------------------------------
	// b, err := json.MarshalIndent(periods, "", "    ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(b))

	// TODO: also grab the point forecast data for the specific city name
	office, err := API.GetOffice()
	if err != nil {
		log.Fatalf("error getting the office name: %v", err)
	}

	// TODO: create interface for card image thumbnails and
	// add those thumbnail images to the scrolling cards

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(800, 450, "Weather Forecast (NOAA)")
	rl.SetTargetFPS(60)

	officePos := rl.NewVector2(
		float32(int32(rl.GetScreenWidth())/2-rl.MeasureText(office, 20)/2),
		20,
	)

	bgTextures := ui.GetBackgroundTextures()

	icons := generateWeatherIcons(periods)
	cards := ui.NewUIScrollingCards(len(periods), 0.15, icons)

	index := 0
	period := periods[index]

	// TODO: refactor bg into UI package as a proper type
	bgKey := "day"
	bgBlendColor := rl.White
	bgBlendDuration := float32(0.5)
	bgBlendTime := float32(0.0)
	bgBlendPercent := float32(1.0)
	bgIsBlending := false

	for !rl.WindowShouldClose() {

		// Handle input events and update the UI ("tick")
		if rl.IsKeyPressed(rl.KeySpace) && !cards.IsScrolling() && !bgIsBlending {
			index = (index + 1) % 14 // index of the selected forecast period.
			period = periods[index]
			bgIsBlending = true
			bgBlendTime = 0.0

			cards.BeginScrolling()
		}

		cards.Update()

		// BEGIN 2D DRAWING ---------------------------------------------------
		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.DarkGray)

			if bgIsBlending {
				bgBlendTime += rl.GetFrameTime()
				bgBlendPercent = bgBlendTime / bgBlendDuration
				if bgBlendTime >= bgBlendDuration {
					bgIsBlending = false
					bgBlendPercent = 1.0
				}
			}

			bgKey, bgBlendColor = ui.GetBackgroundType(period)
			rl.DrawTexture(bgTextures[bgKey], 0, 0, rl.ColorAlpha(bgBlendColor, bgBlendPercent))

			// TODO: add middleground effects / overlays based on weather keywords

			rl.DrawText(office, int32(officePos.X), int32(officePos.Y), 20, rl.RayWhite)

			// TODO: add cross fade on for card thumbnails during animation
			// and update the font to something more readable when pixelated.
			cards.Draw()

			rl.DrawText(
				period.Summary,
				int32(rl.GetScreenWidth()/2)-rl.MeasureText(
					period.Summary, 23)/2,
				80,
				23,
				rl.RayWhite,
			)

			periodName := strings.Split(period.Name, " ")
			if len(periodName) > 0 {
				rl.DrawText(periodName[0], int32(rl.GetScreenWidth()/5)-rl.MeasureText(periodName[0], 30)/2, 140, 30, rl.RayWhite)
			}
			if len(periodName) > 1 {
				rl.DrawText(periodName[1], int32(rl.GetScreenWidth()/5)-rl.MeasureText(periodName[1], 30)/2, 165, 30, rl.RayWhite)
			}

			rl.DrawText(
				fmt.Sprintf("%0.f°%s", period.Temperature, period.TemperatureUnit),
				int32(rl.GetScreenWidth()/2)-rl.MeasureText(
					fmt.Sprintf("%0.f°%s", period.Temperature, period.TemperatureUnit), 100)/2,
				125,
				100,
				rl.RayWhite,
			)

			rl.DrawText(period.WindDirection, int32(rl.GetScreenWidth())-int32(rl.GetScreenWidth()/5)-rl.MeasureText(period.WindDirection, 30)/2, 140, 30, rl.RayWhite)
			rl.DrawText(period.WindSpeed, int32(rl.GetScreenWidth())-int32(rl.GetScreenWidth()/5)-rl.MeasureText(period.WindSpeed, 30)/2, 164, 30, rl.RayWhite)
		}
		rl.EndDrawing()
		// END 2D DRAWING -----------------------------------------------------
	}

	// FREE raylib allocated memory
	// ------------------------------------------------------------------------
	// Internally raylib allocates memory for certain objects (in GPU or RAM)
	// RenderTexture2D, Image, etc. and these should be freed once not needed
	// See: rl.Unload*... methods
	{
		for _, value := range bgTextures {
			rl.UnloadTexture(value)
		}
		for _, value := range icons {
			rl.UnloadTexture(value)
		}
	}

	rl.CloseWindow()
}

// TODO: Draw additional icons / graphics based on weather key words
// and refactor weather icons into the UI package as a proper type.
func generateWeatherIcons(periods []api.ForecastPeriod) []rl.Texture2D {
	icons := make([]rl.Texture2D, len(periods))
	texture := rl.LoadRenderTexture(100, 150)
	size := int32(17)
	for i, period := range periods {
		rl.BeginTextureMode(texture)
		{
			rl.ClearBackground(rl.Blank)
			cardTitle := strings.Split(period.Name, " ")
			switch len(cardTitle) {
			case 1:
				rl.DrawText(cardTitle[0], 50-rl.MeasureText(cardTitle[0], size)/2, 10, size, rl.RayWhite)
			case 2:
				rl.DrawText(cardTitle[0], 50-rl.MeasureText(cardTitle[0], size)/2, 10, size, rl.RayWhite)
				rl.DrawText(cardTitle[1], 50-rl.MeasureText(cardTitle[1], size)/2, 25, size, rl.RayWhite)
			default:
				rl.DrawText("Undefined", 50-rl.MeasureText("Undefined", size)/2, 5, size, rl.RayWhite)
			}
			rl.DrawText(
				fmt.Sprintf("%0.f°%s", period.Temperature, period.TemperatureUnit),
				50-rl.MeasureText(
					fmt.Sprintf("%0.f°%s", period.Temperature, period.TemperatureUnit), size+10)/2,
				60,
				27,
				rl.RayWhite,
			)
		}
		rl.EndTextureMode()
		tempImage := rl.LoadImageFromTexture(texture.Texture)
		icons[(i+6)%14] = rl.LoadTextureFromImage(tempImage) // save generated image as Texture2D for later use
		rl.UnloadImage(tempImage)
	}
	rl.UnloadRenderTexture(texture)
	return icons
}
