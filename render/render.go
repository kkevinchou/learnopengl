package render

import (
	"fmt"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"github.com/kkevinchou/learnopengl/glhelpers"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 800
	height = 600
)

const vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;

void main()
{
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
` + "\x00"

const fragmentShaderSource = `
#version 330 core
out vec4 FragColor;

void main()
{
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
` + "\x00"

var window *sdl.Window

func Setup() {
	window = setupWindow()
	setupGL(window)

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices[0]))*len(vertices), unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	glSrcs, freeFn := gl.Strs(vertexShaderSource)
	defer freeFn()
	gl.ShaderSource(vertexShader, 1, glSrcs, nil)
	gl.CompileShader(vertexShader)

	err := glhelpers.GetGLError(vertexShader, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "SHADER::COMPILE_FAILURE::")
	if err != nil {
		panic(err)
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	glSrcs, freeFn = gl.Strs(fragmentShaderSource)
	defer freeFn()
	gl.ShaderSource(fragmentShader, 1, glSrcs, nil)
	gl.CompileShader(fragmentShader)

	err = glhelpers.GetGLError(fragmentShader, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "SHADER::COMPILE_FAILURE::")
	if err != nil {
		panic(err)
	}

	shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)
}

var vao uint32
var shaderProgram uint32

func Render(delta time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(shaderProgram)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
	window.GLSwap()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
func init() {
	runtime.LockOSThread()
}

func setupGL(window *sdl.Window) *sdl.Window {
	_, err := window.GLCreateContext()
	if err != nil {
		panic(fmt.Sprintf("Failed to create context", err))
	}

	if err = gl.Init(); err != nil {
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

	return window
}

func setupWindow() *sdl.Window {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Sprintf("Failed to init SDL", err))
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(fmt.Sprintf("Failed to create window", err))
	}

	return window
}
