package render

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderProgram struct {
	id uint32
}

func NewShaderProgram(vertexSource []byte, fragmentSource []byte) (ShaderProgram, error) {
	vertexShader, err := compileShader(string(vertexSource), gl.VERTEX_SHADER)
	if err != nil {
		return ShaderProgram{}, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(string(fragmentSource), gl.FRAGMENT_SHADER)
	if err != nil {
		return ShaderProgram{}, err
	}
	defer gl.DeleteShader(fragmentShader)

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

		return ShaderProgram{}, fmt.Errorf("failed to link program: %v", log)
	}

	return ShaderProgram{
		id: program,
	}, nil
}

func (p *ShaderProgram) Use() {
	gl.UseProgram(p.id)
}

func (p *ShaderProgram) BindVertexAttributeData(name string, size int, stride int, offset int) {
	attribute := uint32(gl.GetAttribLocation(p.id, gl.Str(name+"\x00")))
	gl.EnableVertexAttribArray(attribute)
	gl.VertexAttribPointerWithOffset(attribute, int32(size), gl.FLOAT, false, int32(stride), uintptr(offset))
}

func compileShader(shaderSource string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(shaderSource)
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

		return 0, fmt.Errorf("failed to compile %v: %v", shaderSource, log)
	}

	return shader, nil
}
