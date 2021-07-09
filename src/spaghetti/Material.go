package spaghetti

import (
	n "github.com/lachee/noodle"
)

type SliceRender struct {
	shader *n.Shader

	buffVertex n.WebGLBuffer
	buffIndex  n.WebGLBuffer

	a_Vertex n.WebGLAttributeLocation
	a_UV     n.WebGLAttributeLocation
	a_Size   n.WebGLAttributeLocation
	a_Offset n.WebGLAttributeLocation

	u_Sampler    n.WebGLUniformLocation
	u_Window     n.WebGLUniformLocation
	u_Dimension  n.WebGLUniformLocation
	u_Projection n.WebGLUniformLocation
	u_Tiling     n.WebGLUniformLocation

	window   SliceWindow
	vertices []float32
	indicies []uint16
	count    int
}

func NewSliceMaterial() (*SliceRender, error) {
	mat := &SliceRender{}

	shader, shaderError := LoadResourceShader("resource://shader/slice.glsl")
	mat.shader = shader
	if shaderError != nil {
		n.Error("Failed to load the shader", shaderError)
		return nil, shaderError
	}

	mat.buffVertex = n.GL.CreateBuffer()
	mat.buffIndex = n.GL.CreateBuffer()

	mat.a_Vertex = shader.GetAttribLocation("a_Vertex")
	mat.a_UV = shader.GetAttribLocation("a_UV")
	mat.a_Size = shader.GetAttribLocation("a_Size")
	mat.a_Offset = shader.GetAttribLocation("a_Offset")

	mat.u_Sampler = shader.GetUniformLocation("u_Sampler")
	mat.u_Window = shader.GetUniformLocation("u_Window")
	mat.u_Dimension = shader.GetUniformLocation("u_Dimension")
	mat.u_Projection = shader.GetUniformLocation("u_Projection")
	mat.u_Tiling = shader.GetUniformLocation("u_Tiling")

	return mat, nil
}

//beginRender setups GL to use the current state of the material
func (render *SliceRender) setup() {

	var GL = n.GL
	GL.UseProgram(render.shader.GetProgram())

	render.window.texture.SetSampler(render.u_Sampler, 0)

	// Set the vertex data
	stride := 12 + 8 + 8 + 8
	GL.BindBuffer(n.GlArrayBuffer, render.buffVertex)
	GL.VertexAttribPointer(render.a_Vertex, 3, n.GlFloat, false, stride, 0)
	GL.VertexAttribPointer(render.a_UV, 2, n.GlFloat, false, stride, 12)
	GL.VertexAttribPointer(render.a_Size, 2, n.GlFloat, false, stride, 20)
	GL.VertexAttribPointer(render.a_Offset, 2, n.GlFloat, false, stride, 28)
	GL.EnableVertexAttribArray(render.a_Vertex)
	GL.EnableVertexAttribArray(render.a_UV)
	GL.EnableVertexAttribArray(render.a_Size)
	GL.EnableVertexAttribArray(render.a_Offset)

	// Set the projection
	projection := getProjection()
	GL.UniformMatrix4fv(render.u_Projection, projection)

	//TOP LEFT BOTTOM RIGHT
	GL.Uniform4v(render.u_Window, render.window.Window())
	GL.Uniform2v(render.u_Dimension, render.window.Size())
	GL.Uniform2f(render.u_Tiling, float32(render.window.cols), float32(render.window.rows))
}

//SetWindow sets the current cut content
func (render *SliceRender) SetWindow(window SliceWindow) {
	render.window = window
}

//Draw pushes a rectangle into the stack
func (render *SliceRender) Draw(rectangle Rectangle, tile Point, window SliceWindow) {

	// Draw the texture
	if window.texture != render.window.texture {
		// Render the previous content
		render.Render()
	}

	// Update the window
	render.window = window

	size := rectangle.Size()

	// Push Verticies
	render.vertices = append(render.vertices,
		// Vertex, UV, Size, Offset
		rectangle.X, rectangle.Y, 0, // 0 0 0
		0, 0,
		size.X, size.Y,
		float32(tile.X), float32(tile.Y),

		rectangle.X+rectangle.Width, rectangle.Y, 0, // 1 0 0
		1, 0,
		size.X, size.Y,
		float32(tile.X), float32(tile.Y),

		rectangle.X, rectangle.Y+rectangle.Height, 0, // 0 1 0
		0, 1,
		size.X, size.Y,
		float32(tile.X), float32(tile.Y),

		rectangle.X+rectangle.Width, rectangle.Y+rectangle.Height, 0, // 1 1 0
		1, 1,
		size.X, size.Y,
		float32(tile.X), float32(tile.Y),
	)

	// Push indicies
	index := uint16(render.count * 4)
	render.indicies = append(render.indicies,
		index+0, index+1, index+2,
		index+2, index+1, index+3,
	)

	// Increment the count
	render.count++
}

// Render flushes the current arrays and draws the elements. This will setup the materials.
func (render *SliceRender) Render() int {
	if render.count == 0 {
		return 0
	}

	// Setup the material
	render.setup()

	// Draw the buffers
	n.GL.BindBuffer(n.GlArrayBuffer, render.buffVertex)
	n.GL.BufferData(n.GlArrayBuffer, render.vertices, n.GlStaticDraw)

	n.GL.BindBuffer(n.GlElementArrayBuffer, render.buffIndex)
	n.GL.BufferData(n.GlElementArrayBuffer, render.indicies, n.GlStaticDraw)

	n.GL.DrawElements(n.GlTriangles, len(render.indicies), n.GlUnsignedShort, 0)

	// Clear the buffers
	render.vertices = render.vertices[:0]
	render.indicies = render.indicies[:0]

	// Return the count while clearing it
	count := render.count
	render.count = 0
	return count
}

/*
//Draws the vertices and indicies.
func (render *SliceRender) Draw(vertices []float32, indicies []uint16) {

	// Bind texture
	//mat.cut.texture.SetSampler(mat.u_Sampler, 0)
	//defer mat.cut.texture.Unbind()
	n.GL.BindBuffer(n.GlArrayBuffer, render.buffVertex)
	n.GL.BufferData(n.GlArrayBuffer, vertices, n.GlStaticDraw)

	n.GL.BindBuffer(n.GlElementArrayBuffer, render.buffIndex)
	n.GL.BufferData(n.GlElementArrayBuffer, indicies, n.GlStaticDraw)

	n.GL.DrawElements(n.GlTriangles, len(indicies), n.GlUnsignedShort, 0)
	checkGLPanic()
}
*/
type SliceWindow struct {
	top, left, bottom, right float32
	cols, rows               int
	texture                  *n.Texture
}

//Size returns the sliced sice of the texture
func (window SliceWindow) Size() Vector2 {
	return n.NewVector2(float32(window.texture.Width()), float32(window.texture.Height()))
}

//Window gets the window
func (window SliceWindow) Window() Vector4 {
	//size := cut.Size()
	return Vector4{
		X: window.top,
		Y: window.left,
		Z: window.bottom,
		W: window.right,
	}
}
