/*
 * @Author: lmio
 * @Date: 2023-02-27 17:33:58
 * @LastEditTime: 2023-03-01 11:01:46
 * @FilePath: /opengl/glutils/Shader.go
 * @Description:着色器类
 */
package glutils

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"log"
	"os"
	"strings"
)

type Shader struct {
	program uint32
}

func (s *Shader) Use() {
	gl.UseProgram(s.program)
}

func (s *Shader) Delete() {
	gl.DeleteProgram(s.program)
}

func (s *Shader) GetUniform(name string) int32 {
	return gl.GetUniformLocation(s.program, gl.Str(name))
}

// 读取着色器文件
func ReadGLSLFile(filename string) string {

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	return string(content) + "\x00"
}

// 新建着色器程序
func NewShader(vertexShaderSource, fragmentShaderSource string) (*Shader, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
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

		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &Shader{program}, nil
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
