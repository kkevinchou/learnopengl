package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/kkevinchou/learnopengl/render"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gameUpdateDelta = 10 * time.Millisecond
)

var (
	fps = float64(60)
)

func main() {
	rand.Seed(time.Now().Unix())

	previousTime := time.Now()
	var renderAccumulator time.Duration

	msPerFrame := time.Duration(1000000.0/fps) * time.Microsecond

	render.Setup()
	for {
		now := time.Now()
		delta := time.Since(previousTime)
		if delta > 250*time.Millisecond {
			delta = 250 * time.Millisecond
		}
		previousTime = now

		renderAccumulator += delta

		if renderAccumulator >= msPerFrame {
			render.Render(msPerFrame)
		}
		for renderAccumulator > msPerFrame {
			renderAccumulator -= msPerFrame
		}

		sdl.PumpEvents()
		if sdl.GetKeyboardState()[sdl.SCANCODE_ESCAPE] > 0 {
			os.Exit(0)
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
	}
}
