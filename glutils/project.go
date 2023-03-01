/*
 * @Author: lmio
 * @Date: 2023-03-01 15:21:39
 * @LastEditTime: 2023-03-01 16:11:23
 * @FilePath: /opengl/glutils/project.go
 * @Description:投影矩阵
 */
package glutils

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Projection struct {
	projectionUniform int32
	mgl32.Mat4
	fovy, aspect, near, far float32
}

func NewProject(projectionUniform int32, fovy float32, aspect float32, near float32, far float32) *Projection {
	projection := mgl32.Perspective(mgl32.DegToRad(fovy), aspect, near, far)
	return &Projection{
		projectionUniform,
		projection,
		fovy,
		aspect,
		near,
		far,
	}
}

func (p *Projection) Reset() {
	p.Mat4 = mgl32.Perspective(mgl32.DegToRad(p.fovy), p.aspect, p.near, p.far)
	gl.UniformMatrix4fv(p.projectionUniform, 1, false, &p.Mat4[0])
}
