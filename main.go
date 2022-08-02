package main

import (
	"runtime"
	"vox/render"

	_ "embed"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//go:embed main.vs
var vsFile []byte

//go:embed main.fs
var fsFile []byte

func init() {
	runtime.LockOSThread()
}

var vertices = []float32{
	-0.5, -0.5, 0.0,
	0.5, -0.5, 0.0,
	0.0, 0.5, 0.0,
}

var colors = []float32{
	1.0, 0.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 0.0, 1.0,
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Vox", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		panic(err)
	}
	gl.Viewport(0, 0, 800, 600)

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	shaderProgram, err := render.NewShaderProgram(vsFile, fsFile)
	if err != nil {
		panic(err)
	}

	shaderProgram.Use()

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	var vertsBuffer uint32
	gl.GenBuffers(1, &vertsBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertsBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	shaderProgram.BindVertexAttributeData("vert", 3, 3*4, 0)

	gl.GenBuffers(2, &vertsBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertsBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4, gl.Ptr(colors), gl.STATIC_DRAW)
	shaderProgram.BindVertexAttributeData("color", 3, 3*4, 0)

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		shaderProgram.Use()
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
