package spaghetti

import n "github.com/lachee/noodle"

type Slice struct {
	shader       *n.Shader
	vertexBuffer n.WebGLBuffer
	indexBuffer  n.WebGLBuffer
	a_Vertex     n.WebGLAttributeLocation
	a_UV         n.WebGLAttributeLocation
	u_Size       n.WebGLUniformLocation
	u_Window     n.WebGLUniformLocation
	u_Dimension  n.WebGLUniformLocation
	u_Projection n.WebGLUniformLocation
	u_Sampler    n.WebGLUniformLocation
	u_Offset     n.WebGLUniformLocation
	u_Tiling     n.WebGLUniformLocation
	cut          *NineSliceCut
}

func createSlice() (*Slice, error) {
	var GL = n.GL
	s := &Slice{}

	shader, shaderError := LoadResourceShader("resource://shader/slice.glsl")
	if shaderError != nil {
		n.Error("Failed to load the shader", shaderError)
		return nil, shaderError
	}

	/*
		//Load Image
		image, imageError := LoadResourceImage("resource://textures/horizontal_atlas.png")
		if imageError != nil {
			n.Error("Failed to load the image", imageError)
			return nil, imageError
		}

		//Set the buffers and get the locations
		texture := image.CreateTexture()
		texture.SetFilter(n.TextureFilterNearest)
		s.cut = &NineSliceCut{
			top: 5, left: 8, bottom: 5, right: 8,
			cols: 2, rows: 1,
			offsetX: 0, offsetY: 0,

			texture: texture,
		}
	*/

	//Load Image
	image, imageError := LoadResourceImage("resource://textures/slice_atlas.png")
	if imageError != nil {
		n.Error("Failed to load the image", imageError)
		return nil, imageError
	}

	//Set the buffers and get the locations
	texture := image.CreateTexture()
	texture.SetFilter(n.TextureFilterNearest)
	s.cut = &NineSliceCut{
		top: 5, left: 5, bottom: 5, right: 5,
		cols: 6, rows: 1,
		offsetX: 0, offsetY: 0,

		texture: texture,
	}

	s.shader = shader
	s.vertexBuffer = GL.CreateBuffer()
	s.indexBuffer = GL.CreateBuffer()
	s.a_Vertex = s.shader.GetAttribLocation("a_Vertex")
	s.a_UV = s.shader.GetAttribLocation("a_UV")
	s.u_Projection = s.shader.GetUniformLocation("u_Projection")
	s.u_Size = s.shader.GetUniformLocation("u_Size")
	s.u_Window = s.shader.GetUniformLocation("u_Window")
	s.u_Dimension = s.shader.GetUniformLocation("u_Dimension")
	s.u_Sampler = s.shader.GetUniformLocation("u_Sampler")
	s.u_Offset = s.shader.GetUniformLocation("u_Offset")
	s.u_Tiling = s.shader.GetUniformLocation("u_Tiling")

	// Return the object
	return s, nil
}

func (s *Slice) Draw(pos Vector2, size Vector2, tile int) {
	var GL = n.GL

	// Set the program
	GL.UseProgram(s.shader.GetProgram())

	// Bind texture
	s.cut.texture.SetSampler(s.u_Sampler, 0)
	defer s.cut.texture.Unbind()

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
	GL.Uniform2v(s.u_Size, size)

	//TOP LEFT BOTTOM RIGHT
	s.cut.offsetX = tile
	GL.Uniform4v(s.u_Window, s.cut.Window())
	GL.Uniform2v(s.u_Dimension, s.cut.Size())
	GL.Uniform2v(s.u_Offset, s.cut.Offset())
	GL.Uniform2f(s.u_Tiling, float32(s.cut.cols), float32(s.cut.rows))

	//GL.Uniform2fv(bg.u_Resolution, []float32{GL.BoundingBox().Width, GL.BoundingBox().Height})

	// Draw the elements
	n.GL.DrawElements(n.GlTriangles, len(indecies), n.GlUnsignedShort, 0)
}

type NineSliceCut struct {
	top, left, bottom, right float32
	texture                  *n.Texture
	cols, rows               int
	offsetX, offsetY         int
}

//Size returns the sliced sice of the texture
func (cut *NineSliceCut) Size() Vector2 {
	return n.NewVector2(float32(cut.texture.Width()), float32(cut.texture.Height()))
}

//Offset gets teh offset of the slice
func (cut *NineSliceCut) Offset() Vector2 {
	return n.NewVector2i(cut.offsetX, cut.offsetY)
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
