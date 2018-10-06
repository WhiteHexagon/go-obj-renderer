package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Material - rgb color, maybe later texture stuff
type Material struct {
	r, g, b float32
}

// Vertex - vertex
type Vertex struct {
	x, y, z float32
}

// Face - face
type Face struct {
	mtl     string
	indices []int
}

type vkey struct {
	vtx Vertex
	mtl string
}

func localPath(filename string) string {
	p := filepath.Join("objs", filename)
	objpath, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}
	return objpath
}

func threefloatsfrom(arr []string) []float32 {
	a, err := strconv.ParseFloat(arr[0], 32)
	if err != nil {
		panic(err)
	}
	b, err := strconv.ParseFloat(arr[1], 32)
	if err != nil {
		panic(err)
	}
	c, err := strconv.ParseFloat(arr[2], 32)
	if err != nil {
		panic(err)
	}
	result := []float32{float32(a), float32(b), float32(c)}
	return result
}

func loadMaterials(filename string) {
	//fmt.Printf("LoadMaterials=[%v]\n", filename)
	materials = make(map[string]Material)

	file, err := os.Open(localPath(filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	matname := ""
	matcount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "newmtl") {
			arr := strings.Split(line, " ")
			matname = arr[1]
			matcount++
		}

		if strings.HasPrefix(line, "Kd") {
			//TODO this doesnt work, so we use hard coded materials just for now
			//arr := strings.Split(line, " ")
			//mat := threefloatsfrom(arr[1:])
			//materials[matname] = Material{mat[0], mat[1], mat[2]}
			if matname == "grass" {
				materials[matname] = Material{0, 140.0 / 255, 20.0 / 255}
			} else if matname == "earth" {
				materials[matname] = Material{78.0 / 255, 32.0 / 255, 20.0 / 255}
			} else {
				materials[matname] = Material{.5, .5, .5}
			}
		}
	}
}

//face looks like 'f 4//1 6//1 8//1'
func decodeFace(arr []string) (int, []int) {
	var face int
	indices := make([]int, len(arr))
	for i, corner := range arr {
		bit := strings.Split(corner, "/")
		f, err := strconv.ParseInt(bit[2], 10, 32)
		face = int(f)
		if err != nil {
			panic(err)
		}

		index, err := strconv.ParseInt(bit[0], 10, 32)
		indices[i] = int(index) - 1 //need this zero based
		if err != nil {
			panic(err)
		}
	}
	if len(indices) != 3 {
		panic("not triangles...")
	}
	return face, indices
}

var materials map[string]Material
var vertices map[int]Vertex
var faces map[int]Face

func scanFile(filename string) {
	file, err := os.Open(localPath(filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	vertcount := 0
	facecount := 0
	vertices = make(map[int]Vertex)
	faces = make(map[int]Face)
	var mtlname string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "mtllib") {
			arr := strings.Split(line, " ")
			loadMaterials(arr[1])
		}

		if strings.HasPrefix(line, "v ") {
			arr := strings.Split(line, " ")
			v := threefloatsfrom(arr[1:])
			newVertex := Vertex{v[0], v[1], v[2]}
			vertices[vertcount] = newVertex
			vertcount++
		}

		if strings.HasPrefix(line, "usemtl ") {
			arr := strings.Split(line, " ")
			mtlname = arr[1]
		}

		if strings.HasPrefix(line, "f ") {
			arr := strings.Split(line, " ")
			_, indices := decodeFace(arr[1:])
			faces[facecount] = Face{mtlname, indices}
			facecount++
		}

	}
}

func buildVertexInfo(face Face) []VertexInfo {
	var result []VertexInfo
	mat := materials[face.mtl]
	for i := 0; i < len(face.indices); i++ {
		v := vertices[face.indices[i]]
		result = append(result, VertexInfo{v.x, v.y, v.z, mat.r, mat.g, mat.b, 1.0})
	}
	return result
}

//LoadOBJ - Load .obj Model
func LoadOBJ(filename string) *ModelVBO {
	// Parse the file
	scanFile(filename)

	// Build Model
	var result = ModelVBO{}
	uniqueVI := make(map[VertexInfo]bool)
	indexMapping := make(map[vkey]int)
	result.vertices = []VertexInfo{}
	index := 0
	for _, face := range faces { //unknown order!
		vis := buildVertexInfo(face)
		for _, vi := range vis {
			if !uniqueVI[vi] { //if unique
				uniqueVI[vi] = true
				indexMapping[vkey{Vertex{vi.x, vi.y, vi.z}, face.mtl}] = index
				result.vertices = append(result.vertices, vi)
				index++
			}
		}
	}

	// Populate indices
	result.indices = []uint32{}
	for _, face := range faces { //unknown order!
		for i := 0; i < len(face.indices); i++ {
			vtx := vertices[face.indices[i]]
			newIndex := indexMapping[vkey{vtx, face.mtl}]
			result.indices = append(result.indices, uint32(newIndex))
		}
	}

	result.build()
	return &result
}
