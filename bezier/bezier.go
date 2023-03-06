/*
 * @Author: lmio
 * @Date: 2023-03-06 20:58:55
 * @LastEditTime: 2023-03-06 21:58:45
 * @FilePath: /opengl/bezier/bezier.go
 * @Description:绘制bezier曲线
 */
package main

import (
	"fmt"
	"log"
	"opengl/glutils"
	"runtime"
	"github.com/go-gl/gl/v4.1-core/gl"
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
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// 编译着色器程序
	vertexShaderSource := glutils.ReadGLSLFile("bezier/vertexShader.glsl")
	fragmentShaderSource := glutils.ReadGLSLFile("bezier/fragShader.glsl")
	program, err := glutils.NewShader(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatalln("创建着色器程序失败", err)
	}
	defer program.Delete()
	program.Use()

    // Control points for the quadratic Bezier curve.
    controlPoints := []float32{
        0.0, 0.0,
        0.5, 1.0,
        1.0, 0.0,
    }

    // Convert control points to vertices.
    var vertices []mgl32.Vec2
    for t := float32(0); t <= 1.0; t += 0.01 {
        x := (1.0-t)*(1.0-t)*controlPoints[0] + 2*(1.0-t)*t*controlPoints[2] + t*t*controlPoints[4]
        y := (1.0-t)*(1.0-t)*controlPoints[1] + 2*(1.0-t)*t*controlPoints[3] + t*t*controlPoints[5]
        vertices = append(vertices, mgl32.Vec2{x, y})
    }

	fmt.Print(vertices)

	vertexs := glutils.NewVertexsBy2D(vertices)

	vao := glutils.CreateVertexArrayObject(vertexs)

    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 2*4, 0)

	window.Display(func () {
		gl.ClearColor(1.0, 1.0, 1.0, 0.0)
        gl.Clear(gl.COLOR_BUFFER_BIT)
		
		program.Use()
		gl.BindVertexArray(vao)
        gl.DrawArrays(gl.LINE, 0, int32(len(vertices)/2))
		gl.BindVertexArray(0)
	})
}