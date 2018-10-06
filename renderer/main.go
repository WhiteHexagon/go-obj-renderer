package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const screenWidth int = 1280
const screenHeight int = 960

func init() {
	runtime.LockOSThread()
}

func main() {
	render(LoadOBJ("earth_1111.obj"))
}

func render(model *ModelVBO) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(screenWidth, screenHeight, "WhiteHexagon .obj Renderer", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	fmt.Println("OpenGL: ", gl.GoStr(gl.GetString(gl.VERSION)))

	program, err := CreateShaderProgram()
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)
	gl.BindFragDataLocation(program, 0, gl.Str("out_color\x00"))

	// OpenGL setup
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	defer gl.BindVertexArray(0)

	//vertices info
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, model.verticesByteLength, gl.Ptr(model.vertices), gl.STATIC_DRAW)

	//indices
	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, model.indicesByteLength, gl.Ptr(model.indices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("in_position\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, model.verticesStride, gl.PtrOffset(0))

	vertAttrib = uint32(gl.GetAttribLocation(program, gl.Str("in_vertexColor\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 4, gl.FLOAT, false, model.verticesStride, gl.PtrOffset(model.verticesColorOffset))

	gl.Enable(gl.DEPTH_TEST)

	angle := 0.0
	previousTime := glfw.GetTime()

	color := GLColorFromHex("ccffffff")
	gl.ClearColor(color[0], color[1], color[2], color[3])

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearDepth(1.0)

		gl.UseProgram(program)

		projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(screenWidth)/float32(screenHeight), 1.0, 10000.0)
		projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

		camera := mgl32.LookAtV(mgl32.Vec3{-16, 16, -16}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
		cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		modelM := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelUniform, 1, false, &modelM[0])

		gl.BindVertexArray(vao)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
		gl.DrawElements(gl.TRIANGLES, model.indicesCount, gl.UNSIGNED_INT, gl.PtrOffset(0))

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
