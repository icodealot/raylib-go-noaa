package ui

import (
	"fmt"
	"noaawc/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UICard struct {
	Position       rl.Vector2
	TargetPosition rl.Vector2
	Size           rl.Vector2
	TargetSize     rl.Vector2
	Index          int
}

type UIScrollingCards struct {
	Cards     []UICard
	Duration  float32
	animating bool
	runtime   float32
}

func NewUIScrollingCards(count int, duration float32) *UIScrollingCards {
	if count <= 0 || duration < 0 {
		return nil
	}
	// Create a slice of UICards
	cards := make([]UICard, count)
	// Create a UICard for each period
	for i := range cards {
		// Create a new UICard
		scale := util.SinScale(i, len(cards))

		// TODO: make UICard sizes configurable
		width := 100 * scale
		xoffset := (100-width)*0.5 - 335 // 14 cards but only 7 on screen at a time shift left to center
		height := 150 * scale
		yoffset := (150 - height) * 0.5
		padding := 100 / 8 // 7 cards on screen + 1 for the final space

		cards[i] = UICard{
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
	}

	return &UIScrollingCards{
		Cards:     cards,
		Duration:  duration,
		animating: false,
		runtime:   0,
	}
}

func (u *UIScrollingCards) Draw() {
	for _, card := range u.Cards {

		// TODO: calculate the visible range based on configured sizes
		if card.Index < 2 || card.Index > 9 {
			continue
		}

		// half sine wave to ease scale in and out
		//scale := float32(math.Sin((float64(card.Index) - 0.5) / float64(len(cards)) * 180.0))
		scale := util.SinScale(card.Index, len(u.Cards))
		if scale < 0.5 {
			scale = 0.5
		}
		r := rl.Rectangle{
			X:      card.Position.X,
			Y:      card.Position.Y,
			Width:  card.Size.X,
			Height: card.Size.Y,
		}

		rl.DrawRectangleRounded(r, 0.25, 15, rl.ColorAlpha(rl.LightGray, scale))
		rl.DrawRectangleRoundedLines(r, 0.25, 15, 3*scale, rl.ColorAlpha(rl.RayWhite, scale))
		rl.DrawText(fmt.Sprintf("%d", card.Index+1), int32(card.Position.X+card.Size.X/2), int32(card.Position.Y+card.Size.Y/2), 14, rl.RayWhite)
		//fmt.Printf("Card %d, Alpha %2.2f\n", i, scale)
	}
}

func (u *UIScrollingCards) BeginScrolling() {
	if len(u.Cards) > 0 && !u.IsScrolling() {
		// Set new target position for each card.
		u.animating = true
		u.runtime = 0.0
		tempCard := u.Cards[len(u.Cards)-1] // cache the final card for later
		for i := range u.Cards {
			u.Cards[i].Index--
			if u.Cards[i].Index < 0 {
				u.Cards[i].Index = len(u.Cards) - 1
			}

			if i > 0 {
				u.Cards[i].TargetPosition = u.Cards[i-1].Position
				u.Cards[i].TargetSize = u.Cards[i-1].Size
			}
		}
		u.Cards[0].TargetPosition = tempCard.Position
		u.Cards[0].TargetSize = tempCard.Size
		u.Cards[0].Index = tempCard.Index
	}
}

func (u *UIScrollingCards) Update() {
	if u.animating {
		u.runtime += rl.GetFrameTime()
		if u.runtime >= u.Duration {
			// Stop animation and set the cards to their final position
			u.animating = false
			for i := range u.Cards {
				u.Cards[i].Position = u.Cards[i].TargetPosition
				u.Cards[i].Size = u.Cards[i].TargetSize
			}
		} else {
			for i, card := range u.Cards {
				u.Cards[i].Position = rl.Vector2Lerp(card.Position, card.TargetPosition, u.runtime/u.Duration)
				u.Cards[i].Size = rl.Vector2Lerp(card.Size, card.TargetSize, u.runtime/u.Duration)
			}
		}
	}
}

func (u *UIScrollingCards) IsScrolling() bool {
	return u.animating
}
