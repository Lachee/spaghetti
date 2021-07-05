package spaghetti

import n "github.com/lachee/noodle"

type Slice struct {
	shader       *n.Shader
	texture      *n.Texture
	vertexBuffer n.WebGLBuffer
	indexBuffer  n.WebGLBuffer
	a_Vertex     n.WebGLAttributeLocation
	a_UV         n.WebGLAttributeLocation
	u_Projection n.WebGLUniformLocation
	u_Sampler    n.WebGLUniformLocation
}

func createSlice() (*Slice, error) {
	var GL = n.GL
	s := &Slice{}

	shader, err := n.LoadShaderFromCombinedURL("resources/shader/slice.glsl")
	if err != nil {
		return nil, err
	}

	image, err := n.LoadImage("resources/textures/slice.png")
	if err != nil {
		return nil, err
	}

	//Set the buffers and get the locations
	s.texture = image.CreateTexture()
	s.texture.SetFilter(n.TextureFilterNearest)
	s.shader = shader
	s.vertexBuffer = GL.CreateBuffer()
	s.indexBuffer = GL.CreateBuffer()
	s.a_Vertex = s.shader.GetAttribLocation("a_Vertex")
	s.a_UV = s.shader.GetAttribLocation("a_UV")
	s.u_Projection = s.shader.GetUniformLocation("u_Projection")
	s.u_Sampler = s.shader.GetUniformLocation("u_Sampler")

	// Return the object
	return s, nil
}

func (s *Slice) Draw(pos Vector2, size Vector2) {
	var GL = n.GL

	// Bind texture
	s.texture.SetSampler(s.u_Sampler, 0)
	defer s.texture.Unbind()

	mesh := []Vector3{
		n.NewVector3(0, 0, 0).Add(pos.ToVector3()),
		n.NewVector3(size.X, 0, 0).Add(pos.ToVector3()),
		n.NewVector3(0, size.Y, 0).Add(pos.ToVector3()),
		n.NewVector3(size.X, size.Y, 0).Add(pos.ToVector3()),
	}

	uv := []Vector2{
		n.NewVector2(0, 0),
		n.NewVector2(1, 0),
		n.NewVector2(0, 1),
		n.NewVector2(1, 1),
	}

	indecies := []uint16{
		0, 1, 2,
		2, 1, 3,
	}

	var buffer []float32
	for i := 0; i < len(mesh); i++ {
		buffer = append(buffer, mesh[i].Decompose()...)
		buffer = append(buffer, uv[i].Decompose()...)

	}

	// Bind the vertex and Indecides data
	GL.BindBuffer(n.GlArrayBuffer, s.vertexBuffer)
	GL.BufferData(n.GlArrayBuffer, buffer, n.GlStaticDraw)

	// Bind the indicies
	n.GL.BindBuffer(n.GlElementArrayBuffer, s.indexBuffer)
	n.GL.BufferData(n.GlElementArrayBuffer, indecies, n.GlStaticDraw)

	// Use our program
	GL.UseProgram(s.shader.GetProgram())

	// Set the vertex data
	stride := (3 * 4) + (2 * 4)
	GL.BindBuffer(n.GlArrayBuffer, s.vertexBuffer)
	GL.VertexAttribPointer(s.a_Vertex, 3, n.GlFloat, false, stride, 0)
	GL.VertexAttribPointer(s.a_UV, 2, n.GlFloat, false, stride, 3*4)
	GL.EnableVertexAttribArray(s.a_Vertex)
	GL.EnableVertexAttribArray(s.a_UV)

	// Set the projection
	projection := getProjection()
	GL.UniformMatrix4fv(s.u_Projection, projection)
	//GL.Uniform2fv(bg.u_Resolution, []float32{GL.BoundingBox().Width, GL.BoundingBox().Height})

	// Draw the elements
	n.GL.DrawElements(n.GlTriangles, len(indecies), n.GlUnsignedShort, 0)
}
