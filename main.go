package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/learnopengl/render"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width           = 800
	height          = 600
	gameUpdateDelta = 10 * time.Millisecond
)

var (
	fps    = float64(60)
	window *sdl.Window
)

func main() {
	fmt.Println("MAIN")
	rand.Seed(time.Now().Unix())

	// render.Setup()
	previousTime := time.Now()
	var renderAccumulator time.Duration

	msPerFrame := time.Duration(1000000.0/fps) * time.Microsecond

	for {
		now := time.Now()
		delta := time.Since(previousTime)
		if delta > 250*time.Millisecond {
			delta = 250 * time.Millisecond
		}
		previousTime = now

		renderAccumulator += delta

		if renderAccumulator >= msPerFrame {
			render.Render(window, msPerFrame)
		}
		for renderAccumulator > msPerFrame {
			renderAccumulator -= msPerFrame
		}

		sdl.PumpEvents()
		if sdl.GetKeyboardState()[sdl.SCANCODE_ESCAPE] > 0 {
			os.Exit(0)
		}
	}
}

func init() {
	runtime.LockOSThread()

	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Sprintf("Failed to init SDL", err))
	}

	window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(fmt.Sprintf("Failed to create window", err))
	}

	_, err = window.GLCreateContext()
	if err != nil {
		panic(fmt.Sprintf("Failed to create context", err))
	}

	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("Failed to init OpenGL %s", err))
	}

	sdl.SetRelativeMouseMode(false)
	sdl.GLSetSwapInterval(1)

	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.COLOR_MATERIAL)
	gl.ColorMaterial(gl.FRONT, gl.AMBIENT_AND_DIFFUSE)

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
}
