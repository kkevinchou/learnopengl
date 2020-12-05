package render

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var vertices []float32

const vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;

void main()
{
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
`

func Setup() {
	vertices = []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices)), unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	glVertexSource := gl.Str(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, &glVertexSource, nil)
	gl.CompileShader(vertexShader)

	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)

	if status == 0 {
		var errorMessage string
		glErrorMessage := gl.Str(errorMessage)
		gl.GetShaderInfoLog(vertexShader, 512, nil, glErrorMessage)
		fmt.Println(gl.GoStr(glErrorMessage))
	}
}

func Render(window *sdl.Window, delta time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	window.GLSwap()
}
