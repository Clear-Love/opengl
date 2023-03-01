/*
 * @Author: lmio
 * @Date: 2023-02-28 21:47:36
 * @LastEditTime: 2023-03-01 16:32:51
 * @FilePath: /opengl/glutils/camera.go
 * @Description:摄像头
 */
package glutils

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	mgl32.Mat4
	cameraUniform      int32
	pos, front, up     mgl32.Vec3
	speed, sensitivity float32
}

func (c *Camera) Reset() {
	c.Mat4 = mgl32.LookAtV(c.pos, c.pos.Add(c.front), c.up)
	gl.UniformMatrix4fv(c.cameraUniform, 1, false, &c.Mat4[0])
}

func (c *Camera) GetFront() mgl32.Vec3 {
	return c.front
}

func NewCamera(cameraUniform int32, cameraPos, cameraFront, cameraUp mgl32.Vec3) *Camera {
	camera := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
	return &Camera{
		camera,
		cameraUniform,
		cameraPos,
		cameraFront,
		cameraUp,
		0.05,
		0.1,
	}
}
