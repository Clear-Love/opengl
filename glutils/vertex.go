/*
 * @Author: lmio
 * @Date: 2023-03-02 16:31:18
 * @LastEditTime: 2023-03-02 23:26:18
 * @FilePath: /opengl/glutils/vertex.go
 * @Description:顶点结构
 */
package glutils

import "github.com/go-gl/mathgl/mgl32"

// 位置属性
type Position struct {
    mgl32.Vec3	
}

// 法线属性
type Normal struct {
	mgl32.Vec3 
}

// 纹理坐标属性
type TexCoord struct {
	mgl32.Vec2 
} 

// 颜色属性
type Color struct {
	mgl32.Vec3 
}  

// 切线属性
type Tangent struct {
	mgl32.Vec3
}

// 双切线属性
type Bitangent struct {
	mgl32.Vec3 
}

// 三角形面
type Indices [3]uint32 

// 四边形面
type Tetrahedras [4]uint32 
type Vertexs struct {
	vertices []float32  // 顶点数组
	numAttr	 int 		// 属性个数
	attrLen int			// 单个顶点长度
}

func NewVertexs(position []Position) *Vertexs {
	vertices := make([]float32, 0, len(position)*3)
	for _, val := range position {
		vertices = append(vertices, val.X(), val.Y(), val.Z())
	}
	return &Vertexs{vertices, 1, 3}
}

func (v *Vertexs) Addattributev3(attribute []mgl32.Vec3) {
	vertices := make([]float32, 0, len(v.vertices) + len(attribute)*3)
	index := 0
	for _, val := range attribute {
		vertices = append(vertices, v.vertices[index:index+v.attrLen]...)
		vertices = append(vertices, val.X(), val.Y(), val.Z())
		index += v.attrLen + 3
	}
	v.vertices = vertices
}

func (v *Vertexs) Addattributev2(attribute []mgl32.Vec2) {
	nums := make([]float32, 0, len(v.vertices) + len(attribute)*2)
	index := 0
	for _, val := range attribute {
		nums = append(nums, v.vertices[index:index+v.attrLen]...)
		nums = append(nums, val.X(), val.Y())
		index += v.attrLen + 3
	}
	v.vertices = nums
}

func (v *Vertexs) NewNormal(positions []Position, indices []Indices) []Normal {
    normals := make([]Normal, len(positions))
	// 遍历所有三角形，计算每个三角形的法线向量，并将其添加到对应的顶点的法线向量上
	for i := 0; i < len(indices); i++ {
		v0 := positions[indices[i][0]]
		v1 := positions[indices[i+1][1]]
		v2 := positions[indices[i+2][2]]
		
		// 通过叉乘计算法向量
		e1 := v1.Sub(v0.Vec3)
		e2 := v2.Sub(v0.Vec3)
		normal := e1.Cross(e2).Normalize()
		
		// 添加法向量
		normals[i].Vec3, normals[i+1].Vec3, normals[i+2].Vec3  = normal, normal, normal
	}
    return normals
}