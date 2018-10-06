package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/go-gl/gl/all-core/gl"
)

var vertexShaderSource = `
#version 410
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 in_position;
in vec4 in_vertexColor;
out vec4 out_vertexColor;
void main() {
	gl_Position = projection * camera * model * vec4(in_position, 1);
	out_vertexColor = in_vertexColor;
}
` + "\x00"

var fragmentShaderSource = `
#version 410
in vec4 out_vertexColor;
out vec4 out_color;
void main() {
	out_color = out_vertexColor;
}
` + "\x00"

// GLColorFromHex - convert 4 componet hex string into 4 floats for opengl
func GLColorFromHex(input string) [4]float32 {
	decoded, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	var color [4]float32
	for i := range decoded {
		color[i] = float32(decoded[i]) / 255
	}
	return color
}

// CreateShaderProgram - create a shader program
func CreateShaderProgram() (uint32, error) {
	vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

// CompileShader - compile shader
func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Compile failed: %v: %v", source, log)
	}

	return shader, nil
}
