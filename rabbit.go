/*
 * @Author: lmio
 * @Date: 2023-02-18 16:43:54
 * @LastEditTime: 2023-02-26 12:04:59
 * @FilePath: /opengl/rabbit.go
 * @Description:兔子模型
 */
package main

import (
	"log"
	"runtime"
	"opengl/glutils"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
	window := glutils.InitGlfw(width, height)
	defer glfw.Terminate()
	glutils.InitOpenGL()

	// 设置opengl特性
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// 线框模式
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// 编译着色器程序
	program, err := glutils.NewProgram(vertexShaderSource, fragmentShaderSource)
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
	vertices, indices, _, err := glutils.ReadOFFFile("bunny10k.off")
	if err != nil {
		log.Fatalln("failed to read off file!", err)
	}

	vao := glutils.CreateVertexArrayObject(vertices, indices)

	//鼠标滚轮回调函数
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

	//记录鼠标位置
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