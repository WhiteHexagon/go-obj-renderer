package main

import (
	"unsafe"
)

// VertexInfo - single vertex with colour
type VertexInfo struct {
	x, y, z    float32
	r, g, b, a float32
}

// ModelVBO - vetex info and indices for complete model
type ModelVBO struct {
	vertices            []VertexInfo
	verticesColorOffset int
	verticesStride      int32
	verticesByteLength  int
	indices             []uint32
	indicesCount        int32
	indicesByteLength   int
}

func (m *ModelVBO) build() {
	m.verticesColorOffset = 3 * int(unsafe.Sizeof(m.vertices[0].x)) // TODO can we calc this based on position of x in struct? or should we add x+y+z...
	m.verticesStride = int32(unsafe.Sizeof(m.vertices[0]))
	m.verticesByteLength = len(m.vertices) * int(m.verticesStride)
	m.indicesCount = int32(len(m.indices))
	m.indicesByteLength = int(m.indicesCount) * int(unsafe.Sizeof(m.indices[0]))
}
