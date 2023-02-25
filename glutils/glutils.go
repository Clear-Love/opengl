/*
 * @Author: lmio
 * @Date: 2023-02-19 17:55:39
 * @LastEditTime: 2023-02-19 18:50:50
 * @FilePath: /opengl/glutils/glutils.go
 * @Description:
 */
package glutils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

// initGlfw 初始化 glfw，返回一个可用的 Window
func InitGlfw(width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}

// initOpenGL 初始化 OpenGL
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func NewProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
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

func CompileShader(source string, shaderType uint32) (uint32, error) {
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

func ReadOFFFile(filename string) ([]float32, []uint32, []uint32, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	// 读取 OFF 文件头
	var header string
	if _, err := fmt.Fscanf(file, "%s\n", &header); err != nil {
		return nil, nil, nil, err
	}
	if header != "OFF" {
		return nil, nil, nil, fmt.Errorf("invalid OFF file format")
	}

	// 读取顶点和面数
	var numVertices, numFaces, numTetrahedra int
	if _, err := fmt.Fscanf(file, "%d %d %d\n", &numVertices, &numFaces, &numTetrahedra); err != nil {
		return nil, nil, nil, err
	}

	// 读取顶点坐标
	vertices := make([]float32, numVertices*3)
	for i := 0; i < numVertices; i++ {
		if _, err := fmt.Fscanf(file, "%f %f %f\n", &vertices[i*3], &vertices[i*3+1], &vertices[i*3+2]); err != nil {
			return nil, nil, nil, err
		}
	}

	// 读取三角形索引
	indices := make([]uint32, numFaces*3)
	for i := 0; i < numFaces; i++ {
		var n int
		if _, err := fmt.Fscanf(file, "%d %d %d %d\n", &n, &indices[i*3], &indices[i*3+1], &indices[i*3+2]); err != nil {
			return nil, nil, nil, err
		}
		if n != 3 {
			return nil, nil, nil, fmt.Errorf("only triangular meshes are supported")
		}
	}
	//读取四边形索引
	tetrahedras := make([]uint32, numTetrahedra*4)
	for i := 0; i < numTetrahedra; i++ {
		var n int
		if _, err := fmt.Fscanf(file, "%d %d %d %d", &n, &tetrahedras[i*3], &tetrahedras[i*3+1], &tetrahedras[i*3+2]); err != nil {
			return nil, nil, nil, err
		}
		if n != 4 {
			return nil, nil, nil, fmt.Errorf("only tetrahedra meshes are supported")
		}
	}

	return vertices, indices, tetrahedras, nil
}

func CreateVertexArrayObject(vertices []float32, indices []uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// 定义顶点属性指针
	var stride int32 = 3 * 4 // 每个顶点数据占用12个字节（3个float32）
	var offset uintptr = 0
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, offset)
	gl.EnableVertexAttribArray(0)

	return vao
}
