/*
 * @Author: lmio
 * @Date: 2023-02-18 16:43:54
 * @LastEditTime: 2023-03-02 23:49:19
 * @FilePath: /opengl/Rabbit/rabbit.go
 * @Description:兔子模型
 */
package main

import (
	"log"
	"opengl/glutils"
	"runtime"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 1200
	height = 1000
)

func main() {
	runtime.LockOSThread()
	window, err := glutils.NewWindow(width, height)
	if err != nil {
		panic("初始化窗口失败")
	}
	defer window.Terminate()
	glutils.InitOpenGL()

	// 设置opengl特性
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// 线框模式
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// 编译着色器程序
	vertexShaderSource := glutils.ReadGLSLFile("Rabbit/vertexShader.glsl")
	fragmentShaderSource := glutils.ReadGLSLFile("Rabbit/fragmentShader.glsl")
	program, err := glutils.NewShader(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatalln("创建着色器程序失败", err)
	}
	defer program.Delete()
	program.Use()

	// 创建投影矩阵
	projectionUniform := program.GetUniform("projection\x00")
	projection := glutils.NewProject(projectionUniform, 45.0, float32(width)/height, 0.1, 100.0)
	
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection.Mat4[0])

	cameraPos, cameraFront, cameraUp := mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 0, -1}, mgl32.Vec3{0, 1, 0}
	
	cameraUniform := program.GetUniform("camera\x00")
	camera := glutils.NewCamera(cameraUniform, cameraPos, cameraFront, cameraUp)

	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera.Mat4[0])
	model := mgl32.Ident4()
	modelUniform := program.GetUniform("model\x00")
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// 设置清屏颜色
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	// 读取off文件
	vertices, indices, _, err := glutils.ReadOFFFile("Rabbit/bunny10k.off")
	if err != nil {
		log.Fatalln("打开off文件失败:", err)
	}

	vao := glutils.CreateVertexArrayObject(vertices)
	glutils.BinfIndices(indices)

	// 绑定顶点属性指针
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 0, 0) // Position
	//gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 4, 3*4) // Normal
	gl.EnableVertexAttribArray(0)
	//gl.EnableVertexAttribArray(1)

	window.EnableScale(projection)
	window.EnableMoveCameraFront(camera)
	window.EnableMoveCameraPos(camera)

	// 渲染循环
	window.Display(func() {
		// 清空颜色缓冲区和深度缓冲区
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		program.Use()
		camera.Reset()
		projection.Reset()

		// 绘制三角形
		gl.BindVertexArray(vao)
		gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(indices)*3), gl.UNSIGNED_INT, 0)
		gl.BindVertexArray(0)

		// 处理窗口事件
		window.SwapBuffers()
		glfw.PollEvents()
	})

}
