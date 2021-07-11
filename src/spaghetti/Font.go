package spaghetti

import (
	"log"
	"syscall/js"
)

type fontMesh struct {
	verts    []float32
	indicies []uint16
}

type Font struct {
	font  js.Value
	cache map[string]fontMesh
}

// Mesh generates the mesh for the given string
func (f *Font) Mesh(str string, size int) ([]float32, []uint16) {

	// Create the cache if it doesn't exist
	if f.cache == nil {
		f.cache = make(map[string]fontMesh)
	}

	// Return teh cache if we have it
	cacheResult, hasCache := f.cache[str]
	if hasCache {
		return cacheResult.verts, cacheResult.indicies
	}

	result := f.font.Call("mesh", str, size)
	verticies := TypedArrayToFloat32Slice(result.Get("verticies"))
	indices := TypedArrayToUint16Slice(result.Get("indices"))

	// vjs := result.Get("verticies")
	// vlength := vjs.Get("length").Int()
	// verticies = make([]float32, vlength)
	// for i := 0; i < vlength; i++ {
	// 	verticies[i] = float32(vjs.Index(i).Float())
	// }

	// Save to the cache and return result
	log.Println("cache miss", str)
	f.cache[str] = fontMesh{verticies, indices}
	return verticies, indices
}
