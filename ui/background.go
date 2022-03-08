package ui

import (
	"noaawc/api"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO: Setup UIBackground, StartBlending, Update and Draw methods

// CAUTION: these Texture2D objects need to be unloaded
// once "WindowShouldClose" is true.
func GetBackgroundTextures() map[string]rl.Texture2D {

	bgMap := make(map[string]rl.Texture2D)

	tempImage := rl.GenImageGradientV(rl.GetScreenWidth(), rl.GetScreenHeight(), rl.SkyBlue, rl.RayWhite)
	bgMap["day"] = rl.LoadTextureFromImage(tempImage)
	rl.UnloadImage(tempImage)

	nightBlue := rl.NewColor(0, 25, 51, 255)

	tempImage = rl.GenImageGradientV(rl.GetScreenWidth(), rl.GetScreenHeight(), nightBlue, rl.DarkGray)
	bgMap["night"] = rl.LoadTextureFromImage(tempImage)
	rl.UnloadImage(tempImage)

	return bgMap
}

func GetBackgroundType(period api.ForecastPeriod) (string, rl.Color) {
	key := ""
	color := rl.White

	if period.IsDaytime {
		key = "day"
	} else {
		key = "night"
	}

	if strings.Contains(strings.ToLower(period.Summary), "partly cloudy") {
		color = rl.LightGray
	} else if strings.Contains(strings.ToLower(period.Summary), "cloudy") {
		color = rl.DarkGray
	}

	return key, color
}
