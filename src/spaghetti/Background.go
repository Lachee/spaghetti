package spaghetti

import n "github.com/lachee/noodle"

type Background struct {
	shader       *n.Shader
	vertexBuffer n.WebGLBuffer
	a_Vertex     n.WebGLAttributeLocation
	u_Projection n.WebGLUniformLocation
	u_Resolution n.WebGLUniformLocation
}

func createBackground() (*Background, error) {
	var GL = n.GL
	bg := &Background{}

	shader, err := n.LoadShaderFromCombinedURL("resources/shader/background.glsl")
	if err != nil {
		return nil, err
	}

	//Set the buffers and get the locations
	bg.shader = shader
	bg.vertexBuffer = GL.CreateBuffer()
	bg.a_Vertex = bg.shader.GetAttribLocation("a_Vertex")
	bg.u_Projection = bg.shader.GetUniformLocation("u_Projection")
	bg.u_Resolution = bg.shader.GetUniformLocation("u_Resolution")

	// Set the vertex data
	vertices := []Vector3{
		n.NewVector3(0, 0, 0),
		n.NewVector3(1, 0, 0),
		n.NewVector3(0, 1, 0),
		n.NewVector3(1, 1, 0),
		n.NewVector3(0, 1, 0),
		n.NewVector3(1, 0, 0),
	}
	GL.BindBuffer(n.GlArrayBuffer, bg.vertexBuffer)
	GL.BufferData(n.GlArrayBuffer, vertices, n.GlStaticDraw)

	// Return the object
	return bg, nil
}

func (bg *Background) Draw() {
	var GL = n.GL

	GL.UseProgram(bg.shader.GetProgram())

	// Set the vertex data
	GL.BindBuffer(n.GlArrayBuffer, bg.vertexBuffer)
	GL.VertexAttribPointer(bg.a_Vertex, 3, n.GlFloat, false, 0, 0)
	GL.EnableVertexAttribArray(bg.a_Vertex)

	// Set the projection
	projection := getProjection()
	GL.UniformMatrix4fv(bg.u_Projection, projection)
	GL.Uniform2fv(bg.u_Resolution, []float32{GL.BoundingBox().Width, GL.BoundingBox().Height})

	// Draw the element
	GL.DrawArrays(n.GlTriangles, 0, 6)
}
