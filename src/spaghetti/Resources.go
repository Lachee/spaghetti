package spaghetti

import (
	"errors"
	"strings"
	"syscall/js"

	n "github.com/lachee/noodle"
)

type ResourceResult struct {
	js    js.Value
	Error error
}

//FetchResource returns a new resource
func FetchResource(url string) <-chan *ResourceResult {
	channel := make(chan *ResourceResult)
	go func() {
		defer close(channel)
		result := <-ResolvePromise(JS.Call("fetchResource", url))
		if result.Error != nil {
			channel <- &ResourceResult{Error: result.Error}
		} else {

			channel <- &ResourceResult{js: result.Values[0]}
		}
	}()
	return channel
}

//ToBytes converts the JS bytes to a GO bytes
func (result *ResourceResult) ToBytes() []byte {
	buffer := make([]byte, result.js.Get("length").Int())
	js.CopyBytesToGo(buffer, result.js)
	return buffer
}

//ToString converts the bytes to a string
func (result *ResourceResult) ToString() string {
	return string(result.ToBytes())
}

//ToImage loads the JS value as an image. This requires the resource to be a Image or ImageData in JS.
func (result *ResourceResult) ToImage() *n.Image {
	return n.LoadImageJS(result.js)
}

//LoadResourceImage fetches the image from the given resource asyncronously
func LoadResourceImage(resource string) (*n.Image, error) {
	imageResource := <-FetchResource(resource)
	if imageResource.Error != nil {
		return nil, imageResource.Error
	}
	return imageResource.ToImage(), nil
}

//FetchShaderResource fetches the combined shaders from the given resource
func LoadResourceShader(resource string) (*n.Shader, error) {

	shaderResource := <-FetchResource(resource)
	if shaderResource.Error != nil {
		return nil, shaderResource.Error
	}

	// Convert to string
	str := shaderResource.ToString()

	// Prepare the frags
	var fragShaderCode, vertShaderCode string
	indexOfVert := strings.LastIndex(str, "//vert:")
	if indexOfVert < 0 {
		err := errors.New("cannot find vert tag")
		return nil, err
	}
	indexOfFrag := strings.LastIndex(str, "//frag:")
	if indexOfFrag < 0 {
		err := errors.New("cannot find frag tag")
		return nil, err
	}

	if indexOfVert < indexOfFrag {
		vertShaderCode = str[indexOfVert:indexOfFrag]
		fragShaderCode = str[indexOfFrag:]
	} else {
		fragShaderCode = str[indexOfFrag:indexOfVert]
		vertShaderCode = str[indexOfVert:]
	}

	// Load the shader
	return n.LoadShader(vertShaderCode, fragShaderCode)
}
