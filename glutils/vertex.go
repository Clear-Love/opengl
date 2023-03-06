/*
 * @Author: lmio
 * @Date: 2023-03-02 16:31:18
 * @LastEditTime: 2023-03-06 21:24:45
 * @FilePath: /opengl/glutils/vertex.go
 * @Description:顶点结构
 */
package glutils

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// 三角形面
type Indices [3]uint32 

// 四边形面
type Tetrahedras [4]uint32

type Vertexs struct {
	vertices []float32  // 顶点数组
	numAttr	 int 		// 属性个数
	attrLen int			// 单个顶点长度
	attNum int 			// 顶点个数
}

func NewVertexsBy3D(position []mgl32.Vec3) *Vertexs {
	vertices := make([]float32, 0, len(position)*3)
	for _, val := range position {
		vertices = append(vertices, val.X(), val.Y(), val.Z())
	}
	
	return &Vertexs{vertices, 1, 3, len(position)}
}

func NewVertexsBy2D(position []mgl32.Vec2) *Vertexs {
	vertices := make([]float32, 0, len(position)*2)
	for _, val := range position {
		vertices = append(vertices, val.X(), val.Y())
	}
	
	return &Vertexs{vertices, 1, 2, len(position)}
}

func (v *Vertexs) Addattributev3(attribute []mgl32.Vec3) error {
	if v.attNum < len(attribute) {
		return fmt.Errorf("属性数组长度出错")
	}
	vertices := make([]float32, 0, len(v.vertices) + len(attribute)*3)
	start := 0
	for i := 0; i < v.attNum; i++ {
		vertices = append(vertices, v.vertices[start:start+v.attrLen]...)
		vertices = append(vertices, attribute[i][0], attribute[i][1], attribute[i][2])
		start += v.attrLen
	}
	v.vertices = vertices
	v.numAttr++
	v.attrLen += 3
	return nil
}

func (v *Vertexs) Addattributev2(attribute []mgl32.Vec2) error {
	if v.attNum < len(attribute) {
		return fmt.Errorf("属性数组长度出错")
	}
	vertices := make([]float32, 0, len(v.vertices) + len(attribute)*2)
	start := 0
	for i := 0; i < v.attNum; i++ {
		vertices = append(vertices, v.vertices[start:start+v.attrLen]...)
		vertices = append(vertices, attribute[i][0], attribute[i][1])
		start += v.attrLen
	}
	v.vertices = vertices
	v.numAttr++
	v.attrLen += 2
	return nil
}

func NewNormal(positions []mgl32.Vec3, indices []Indices) []mgl32.Vec3 {
    normals := make([]mgl32.Vec3, len(positions))
	// 遍历所有三角形，计算每个三角形的法线向量，并将其添加到对应的顶点的法线向量上
	for i := 0; i < len(indices); i++ {
		ind := indices[i]
		v0 := positions[ind[0]]
		v1 := positions[ind[1]]
		v2 := positions[ind[2]]
		
		// 通过叉乘计算法向量
		e1 := v1.Sub(v0)
		e2 := v2.Sub(v0)
		normal := e1.Cross(e2).Normalize()
		
		// 添加法向量
		normals[ind[0]] = normals[ind[0]].Add(normal)
		normals[ind[1]] = normals[ind[1]].Add(normal)
		normals[ind[2]] = normals[ind[2]].Add(normal)
	}
    return normals
}