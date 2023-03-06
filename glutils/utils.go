/*
 * @Author: lmio
 * @Date: 2023-03-01 00:17:21
 * @LastEditTime: 2023-03-06 21:22:43
 * @FilePath: /opengl/glutils/utils.go
 * @Description:常用函数
 */
package glutils

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// initOpenGL 初始化 OpenGL
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func ReadOFFFile(filename string) ([]mgl32.Vec3, []Indices, []Tetrahedras, error) {
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
	position := make([]mgl32.Vec3, numVertices)
	for i := 0; i < numVertices; i++ {
		var x, y, z float32
		_, err := fmt.Fscanf(file, "%f %f %f\n", &x, &y, &z)
		position[i] = mgl32.Vec3{x, y, z}
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// 读取三角形索引
	indices := make([]Indices, numFaces)
	for i := 0; i < numFaces; i++ {
		var n int
		if _, err := fmt.Fscanf(file, "%d %d %d %d\n", &n, &indices[i][0], &indices[i][1], &indices[i][2]); err != nil {
			return nil, nil, nil, err
		}
		if n != 3 {
			return nil, nil, nil, fmt.Errorf("only triangular meshes are supported")
		}
	}
	//读取四边形索引
	tetrahedras := make([]Tetrahedras, numTetrahedra*4)
	for i := 0; i < numTetrahedra; i++ {
		var n int
		if _, err := fmt.Fscanf(file, "%d %d %d %d %d", &n, &tetrahedras[i][0], &tetrahedras[i][1], &tetrahedras[i][2], &tetrahedras[i][3]); err != nil {
			return nil, nil, nil, err
		}
		if n != 4 {
			return nil, nil, nil, fmt.Errorf("only tetrahedra meshes are supported")
		}
	}

	return position, indices, tetrahedras, nil
}

func BinfIndices(indices []Indices) {
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*12, gl.Ptr(indices), gl.STATIC_DRAW)
}

// 给定顶点坐标，返回顶点数组对象
func CreateVertexArrayObject(v *Vertexs) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(v.vertices)*v.attrLen*4, gl.Ptr(v.vertices), gl.STATIC_DRAW)

	return vao
}