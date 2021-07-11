package spaghetti

import (
	"errors"
	"strings"
	"syscall/js"

	n "github.com/lachee/noodle"
)

/**
ResourceResult is a tuple that contains the JS value and any errors that were created from the resource.
Spaghetti itself contains nothing on the Go side to resolve the resources, that is all handled with the wrapper spaghetti.js module which implements functionality to resolve the resource:// url.
*/
type ResourceResult struct {
	MimeType string
	Data     js.Value
	Error    error
}

func LoadResourceFont(resource string) (*Font, error) {
	// TODO:
	// 1. load the resource
	// 2. invoke the data with the string and font size
	// 3. result of invoke is polygon, use https://github.com/tchayen/triangolatte or https://github.com/mapbox/earcut
	// 4. render the triangulated polygon
	fontResource := <-FetchResource(resource)
	if fontResource.Error != nil {
		return nil, fontResource.Error
	}

	if !strings.HasPrefix(fontResource.MimeType, "font/") {
		return nil, errors.New("Font resource has a invalid MimeType of " + fontResource.MimeType)
	}

	return &Font{font: fontResource.Data}, nil
}

//LoadResourceImage fetches the image from the given resource address. If the resource address is not an image, then an error will be thrown
func LoadResourceImage(resource string) (*n.Image, error) {
	imageResource := <-FetchResource(resource)
	if imageResource.Error != nil {
		return nil, imageResource.Error
	}
	return imageResource.ToImage()
}

//FetchShaderResource fetches the combined shaders from the given resource address.
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

/**
FetchResource returns a new resource.
The urls are passed to the wrapper spaghetti.js module and the promises are resolved.
Resources should be loaded using the resource:// protocol. Spaghetti.js will resolve these paths to either URLs or actual packaged data. Normal urls can still be used in this function.
The data returned is of type Promise<Uint8Array|Image> and ResourceResult provides methods to convert the data to appropriate types.
*/
func FetchResource(url string) <-chan *ResourceResult {
	channel := make(chan *ResourceResult)
	go func() {
		defer close(channel)
		result := <-ResolvePromise(JS.Call("fetchResource", url))
		if result.Error != nil {
			channel <- &ResourceResult{Error: result.Error} // We have an error, abort
		} else {
			arr := result.Values[0]
			channel <- &ResourceResult{MimeType: arr.Index(0).String(), Data: arr.Index(1)} // We succeeded, so lets pipe that good stuff in
		}
	}()
	return channel
}

//ToBytes converts the JS bytes to a GO bytes
func (result *ResourceResult) ToBytes() []byte {
	buffer := make([]byte, result.Data.Get("length").Int())
	js.CopyBytesToGo(buffer, result.Data)
	return buffer
}

//ToString converts the bytes to a string
func (result *ResourceResult) ToString() string {
	if strings.HasPrefix(result.MimeType, "text/") {
		return result.Data.String()
	}
	return string(result.ToBytes())
}

//ToImage loads the JS value as an image. This requires the resource to be a Image or ImageData in JS.
func (result *ResourceResult) ToImage() (*n.Image, error) {
	if !strings.HasPrefix(result.MimeType, "image/") {
		return nil, errors.New("ResourceResult is not an image")
	}

	return n.LoadImageJS(result.Data), nil
}
