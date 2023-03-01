/*
 * @Author: lmio
 * @Date: 2023-02-19 17:55:39
 * @LastEditTime: 2023-03-01 16:31:19
 * @FilePath: /opengl/glutils/window.go
 * @Description:
 */
package glutils

import (
	"math"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	*glfw.Window
	mouseDown  bool
	Posx, Posy float64
}

// initGlfw 初始化 glfw，返回一个可用的 Window
func NewWindow(width, height int) (*Window, error) {
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
		return nil, err
	}
	window.MakeContextCurrent()
	return &Window{window, false, 0, 0}, nil
}

func (w *Window) Display(f func()) {
	for !w.ShouldClose() {
		f()
	}
}

func (w *Window) Terminate() {
	glfw.Terminate()
}

func (w *Window) listenMouse() {
	//鼠标按键回调函数，监听鼠标是否按下
	w.SetMouseButtonCallback(func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				w.Posx, w.Posy = w.GetCursorPos()
				w.mouseDown = true
			} else if action == glfw.Release {
				w.mouseDown = false
			}
		}
	})
}

func (w *Window) EnableMoveCameraFront(c *Camera) {
	w.listenMouse()
	var pitch, yaw float32 = 0, -90

	w.SetCursorPosCallback(func(window *glfw.Window, xpos float64, ypos float64) {
		if !w.mouseDown {
			w.Posx = xpos
			w.Posy = ypos
			return
		}
		xoffset := xpos - w.Posx
		yoffset := w.Posy - ypos
		w.Posx = xpos
		w.Posy = ypos

		// 移动因子
		xoffset *= float64(c.sensitivity)
		yoffset *= float64(c.sensitivity)

		yaw += float32(xoffset)
		pitch += float32(yoffset)

		if pitch > 89 {
			pitch = 89
		}
		if pitch < -89 {
			pitch = -89
		}

		// 计算旋转矩阵
		c.front = mgl32.Vec3{
			float32(math.Cos(float64(mgl32.DegToRad(float32(yaw))) * math.Cos(float64(mgl32.DegToRad(float32(pitch)))))),
			float32(math.Sin(float64(mgl32.DegToRad(float32(pitch))))),
			float32(math.Sin(float64(mgl32.DegToRad(float32(yaw))) * math.Cos(float64(mgl32.DegToRad(float32(pitch)))))),
		}.Normalize()
	})
}

func (w *Window) EnableMoveCameraPos(c *Camera) {
	w.SetKeyCallback(func (window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

		if key == glfw.KeyW && action == glfw.Press {
			c.pos = c.pos.Add(c.front.Mul(c.speed))
		}
		if key == glfw.KeyS && action == glfw.Press {
			c.pos = c.pos.Sub(c.front.Mul(c.speed))
		}
		if key == glfw.KeyA && action == glfw.Press {
			c.pos = c.pos.Sub(((c.front.Cross(c.up)).Mul(c.speed)).Normalize())
		}
		if key == glfw.KeyD && action == glfw.Press {
			c.pos = c.pos.Add(((c.front.Cross(c.up)).Mul(c.speed)).Normalize())
		}
	})
}

func (w *Window) EnableScale(p *Projection) {
	w.SetScrollCallback(func(window *glfw.Window, xoff float64, yoff float64) {
		p.fovy -= float32(yoff)
		if p.fovy < 0 {
			p.fovy = 0
		}
	})
}
