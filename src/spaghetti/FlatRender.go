package spaghetti

import (
	n "github.com/lachee/noodle"
)

type FlatRender struct {
	shader *n.Shader

	buffVertex n.WebGLBuffer
	buffIndex  n.WebGLBuffer

	a_Vertex     n.WebGLAttributeLocation
	a_Color      n.WebGLAttributeLocation
	u_Projection n.WebGLUniformLocation

	vertices []float32
	indicies []uint16
}

func NewFlatRender() (*FlatRender, error) {
	mat := &FlatRender{}

	shader, shaderError := LoadResourceShader("resource://shader/flat.glsl")
	mat.shader = shader
	if shaderError != nil {
		n.Error("Failed to load the flat shader", shaderError)
		return nil, shaderError
	}

	mat.buffVertex = n.GL.CreateBuffer()
	mat.buffIndex = n.GL.CreateBuffer()

	mat.a_Vertex = shader.GetAttribLocation("a_Vertex")
	mat.a_Color = shader.GetAttribLocation("a_Color")
	mat.u_Projection = shader.GetUniformLocation("u_Projection")

	return mat, nil
}

func (render *FlatRender) setup() {
	var GL = n.GL
	GL.UseProgram(render.shader.GetProgram())

	// Set the vertex data
	stride := 12 + 4
	GL.BindBuffer(n.GlArrayBuffer, render.buffVertex)
	GL.VertexAttribPointer(render.a_Vertex, 3, n.GlFloat, false, stride, 0)
	GL.VertexAttribPointer(render.a_Color, 4, n.GlUnsignedByte, true, stride, 12)
	GL.EnableVertexAttribArray(render.a_Vertex)
	GL.EnableVertexAttribArray(render.a_Color)

	// Set the projection
	projection := getProjection()
	GL.UniformMatrix4fv(render.u_Projection, projection)
}

// Draws a rectangle onto the screen
func (render *FlatRender) DrawRectangle(rectangle Rectangle, color Color) {
	tint := color.ToTint()
	startVertex := uint16(len(render.vertices) / 4)

	// Push Verticies
	render.vertices = append(render.vertices,
		// Vertex, UV, Size, Offset
		rectangle.X, rectangle.Y, 0, // 0 0 0
		tint,

		rectangle.X+rectangle.Width, rectangle.Y, 0, // 1 0 0
		tint,

		rectangle.X, rectangle.Y+rectangle.Height, 0, // 0 1 0
		tint,

		rectangle.X+rectangle.Width, rectangle.Y+rectangle.Height, 0, // 1 1 0
		tint,
	)

	// Push indicies
	render.indicies = append(render.indicies,
		startVertex+0, startVertex+1, startVertex+2,
		startVertex+2, startVertex+1, startVertex+3,
	)
}

func (render *FlatRender) DrawText(position Vector2, str string, size int, font *Font, color Color) {
	//+-+waindexOffset := uint16(len(render.vertices) / 4)
	verts, indicies := font.Mesh(str, size)
	tint := color.ToTint()

	// Push Verts
	for i := 0; i < len(verts); i += 2 {
		render.vertices = append(render.vertices,
			// Vertex, UV, Size, Offset
			position.X+verts[i], position.Y+verts[i+1], 0, // 0 0 0
			tint,
		)
	}

	// Push indicies
	for _, i := range indicies {
		render.indicies = append(render.indicies, i)
	}

	/*
		// Push indicies
		for i := 0; i < len(indicies); i += 3 {
			render.indicies = append(render.indicies,
				indexOffset+indicies[i],
				indexOffset+indicies[i+1],
				indexOffset+indicies[i+2],
			)
		}
	*/
}

//Render the buffered mesh
func (render *FlatRender) Render() {

	// Skip if empty
	if len(render.vertices) == 0 {
		return
	}

	// Setup material
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
}
