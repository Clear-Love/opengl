/*
 * @Author: lmio
 * @Date: 2023-03-01 00:17:21
 * @LastEditTime: 2023-03-01 00:23:33
 * @FilePath: /opengl/glutils/utils.go
 * @Description:常用函数
 */
package glutils

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"log"
	"os"
)

// initOpenGL 初始化 OpenGL
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
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
		if _, err := fmt.Fscanf(file, "%d %d %d %d %d", &n, &tetrahedras[i*4], &tetrahedras[i*4+1], &tetrahedras[i*4+2], &tetrahedras[i*4+3]); err != nil {
			return nil, nil, nil, err
		}
		if n != 4 {
			return nil, nil, nil, fmt.Errorf("only tetrahedra meshes are supported")
		}
	}

	return vertices, indices, tetrahedras, nil
}

// 给定顶点坐标和三角形索引，返回顶点数组对象
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
