/*
 * @Author: lmio
 * @Date: 2023-02-18 16:43:54
 * @LastEditTime: 2023-02-20 22:53:29
 * @FilePath: /opengl/rabbit.go
 * @Description:兔子模型
 */
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width              = 1200
	height             = 1000
	vertexShaderSource = `
	#version 410
	uniform mat4 projection;
	uniform mat4 camera;
	uniform mat4 model;
	in vec3 vert;
	in vec2 vertTexCoord;
	out vec4 fragPos;
	void main() {
		fragPos = model * vec4(vert, 1.0);
		gl_Position = projection * camera * model * vec4(vert, 10); 
	}
	` + "\x00"

	fragmentShaderSource = `
	#version 410
	uniform sampler2D tex;
	in vec4 fragPos;
	out vec4 outputColor;
	void main() {
		vec3 color = vec3(0.46, 0.51, 0.64);
		outputColor = vec4(color, 1.0);
	}
	` + "\x00"
)

func main() {
	runtime.LockOSThread()
	window := InitGlfw(width, height)
	defer glfw.Terminate()
	InitOpenGL()

	// 设置opengl特性
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// 线框模式
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// 编译着色器程序
	program, err := NewProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatalln("failed to create program:", err)
	}
	defer gl.DeleteProgram(program)
	gl.UseProgram(program)

	// 创建投影矩阵
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// 设置清屏颜色
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	// 读取off文件
	vertices, indices, _, err := ReadOFFFile("bunny10k.off")
	if err != nil {
		log.Fatalln("failed to read off file!", err)
	}

	vao := CreateVertexArrayObject(vertices, indices)

	// 鼠标滚轮回调函数
	scrollCallback := func(window *glfw.Window, xoff float64, yoff float64) {
		// 计算缩放因子
		scaleFactor := 1.0 + float32(yoff)*0.1

		// 将缩放因子应用到缩放变换矩阵中
		scaleMatrix := mgl32.Scale3D(scaleFactor, scaleFactor, scaleFactor)

		// 更新视图矩阵
		camera = camera.Mul4(scaleMatrix)

		// 将新的投影矩阵传递给着色器
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	}
	window.SetScrollCallback(scrollCallback)

	var mouseDown bool

	//鼠标按键回调函数，监听鼠标是否按下
	mouseButtonCallback := func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				mouseDown = true
			} else if action == glfw.Release {
				mouseDown = false
			}
		}
	}
	window.SetMouseButtonCallback(mouseButtonCallback)

	// 记录鼠标位置
	lastX, lastY := 0.0, 0.0

	cursorPosCallback := func(window *glfw.Window, xpos float64, ypos float64) {
		if !mouseDown {
			lastX = xpos
			lastY = ypos
			return
		}
		xoffset := xpos - lastX
		yoffset := lastY - ypos 
		lastX = xpos
		lastY = ypos

		// 移动因子
		sensitivity := 0.005
		xoffset *= sensitivity
		yoffset *= sensitivity

		camera = mgl32.Translate3D(float32(xoffset), float32(yoffset), 0).Mul4(camera)
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	}
	window.SetCursorPosCallback(cursorPosCallback)

	angle := 0.0
	previousTime := glfw.GetTime()
	// 渲染循环
	for !window.ShouldClose() {
		// 清空颜色缓冲区和深度缓冲区
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		angle += elapsed

		// 计算旋转角度
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// 使用着色器程序
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		// 绘制三角形
		gl.BindVertexArray(vao)
		gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, 0)
		gl.BindVertexArray(0)
		// 处理窗口事件
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

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
	// 启用鼠标滚轮事件
	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetInputMode(glfw.StickyMouseButtonsMode, glfw.True)
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
