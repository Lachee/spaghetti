package spaghetti

import (
	n "github.com/lachee/noodle"
)

type NineSliceMaterial struct {
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

	cut NineSliceCut
}

func NewNineSliceMaterial() (*NineSliceMaterial, error) {
	mat := &NineSliceMaterial{}

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

//UseMaterial setups GL to use the current state of the material
func (mat *NineSliceMaterial) UseMaterial() {

	var GL = n.GL
	GL.UseProgram(mat.shader.GetProgram())

	mat.cut.texture.SetSampler(mat.u_Sampler, 0)

	// Set the vertex data
	stride := 12 + 8 + 8 + 8
	GL.BindBuffer(n.GlArrayBuffer, mat.buffVertex)
	GL.VertexAttribPointer(mat.a_Vertex, 3, n.GlFloat, false, stride, 0)
	GL.VertexAttribPointer(mat.a_UV, 2, n.GlFloat, false, stride, 12)
	GL.VertexAttribPointer(mat.a_Size, 2, n.GlFloat, false, stride, 20)
	GL.VertexAttribPointer(mat.a_Offset, 2, n.GlFloat, false, stride, 28)
	GL.EnableVertexAttribArray(mat.a_Vertex)
	GL.EnableVertexAttribArray(mat.a_UV)
	GL.EnableVertexAttribArray(mat.a_Size)
	GL.EnableVertexAttribArray(mat.a_Offset)

	// Set the projection
	projection := getProjection()
	GL.UniformMatrix4fv(mat.u_Projection, projection)

	//TOP LEFT BOTTOM RIGHT
	GL.Uniform4v(mat.u_Window, mat.cut.Window())
	GL.Uniform2v(mat.u_Dimension, mat.cut.Size())
	GL.Uniform2f(mat.u_Tiling, float32(mat.cut.cols), float32(mat.cut.rows))
	checkGLPanic()
}

//SetCut sets the current cut content
func (mat *NineSliceMaterial) SetCut(cut NineSliceCut) {
	mat.cut = cut
}

//Draws the vertices and indicies.
func (mat *NineSliceMaterial) Draw(vertices []float32, indicies []uint16) {

	// Bind texture
	//mat.cut.texture.SetSampler(mat.u_Sampler, 0)
	//defer mat.cut.texture.Unbind()

	mat.UseMaterial()
	n.GL.BindBuffer(n.GlArrayBuffer, mat.buffVertex)
	n.GL.BufferData(n.GlArrayBuffer, vertices, n.GlStaticDraw)

	n.GL.BindBuffer(n.GlElementArrayBuffer, mat.buffIndex)
	n.GL.BufferData(n.GlElementArrayBuffer, indicies, n.GlStaticDraw)

	n.GL.DrawElements(n.GlTriangles, len(indicies), n.GlUnsignedShort, 0)
	checkGLPanic()
}

func checkGLPanic() {
	err := n.GL.Error()
	if err != nil {
		panic(err)
	}
}

type NineSliceCut struct {
	top, left, bottom, right float32
	texture                  *n.Texture
	cols, rows               int
}

//Size returns the sliced sice of the texture
func (cut *NineSliceCut) Size() Vector2 {
	return n.NewVector2(float32(cut.texture.Width()), float32(cut.texture.Height()))
}

//Window gets the window
func (cut *NineSliceCut) Window() Vector4 {
	//size := cut.Size()
	return Vector4{
		X: cut.top,
		Y: cut.left,
		Z: cut.bottom,
		W: cut.right,
	}
}
