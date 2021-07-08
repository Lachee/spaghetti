package spaghetti

import (
	"log"
	"syscall/js"

	n "github.com/lachee/noodle"
)

// Library is the hook to the JS library
var JS js.Value

var cursorSize = Vector2{1, 1}
var cursorHotspot = Vector2{-1, -3}

//Application handles the game. Put your variables in here
type Application struct {
	shader *n.Shader
	bg     *Background
	slice  *Slice
}

//EnableDebugger will set flags for debugging purposes
func (app *Application) EnableDebugger() {

}

//Start allows for setup
func (app *Application) Start() bool {

	//Load defaults
	JS = n.Canvas().Get("editor")
	pixelScale = 2

	// Create the background
	bg, err := createBackground()
	app.bg = bg
	if err != nil {
		n.Error("Failed to create the background", err)
		return false
	}

	// Create the slice
	s, err := createSlice()
	app.slice = s
	if err != nil {
		n.Error("Failed to create the slice", err)
		return false
	}

	// The mouse should trigger render events
	n.MouseDraws = true
	return true
}

var toggle bool = false

//Update runs once a frame
func (app *Application) Update(dt float32) {
	if n.Input().GetKeyDown(n.KeySpace) {
		toggle = !toggle
		log.Println("toggle", toggle)
	}

	// Pixel Scale change
	if n.Input().GetKey(n.KeyNumAdd) {
		pixelScale += 1
		log.Println("Pixel Scale:", pixelScale)
	} else if n.Input().GetKey(n.KeyNumSubtract) {
		pixelScale -= 1
		if pixelScale <= 0 {
			pixelScale = 1
		}
		log.Println("Pixel Scale:", pixelScale)
	}

	// Camera Movement
	axis := n.Input().GetAxis2D(n.KeyA, n.KeyD, n.KeyW, n.KeyS)
	cameraPosition = cameraPosition.Add(axis.Scale(dt).ToVector3())
	//log.Println("Camera Position", cameraPosition, axis)
}

var (
	indexBuffer  n.WebGLBuffer
	vertexBuffer n.WebGLBuffer

	a_Vertex     n.WebGLAttributeLocation
	u_Projection n.WebGLUniformLocation

	cameraPosition Vector3
	pixelScale     float32
)

//Render draws the frame
func (app *Application) Render() {

	n.DebugDraw = true

	// Clear the canvas
	n.GL.ClearColor(n.White)
	n.GL.Clear(n.GlColorBufferBit | n.GlDepthBufferBit)
	app.bg.Draw()

	var width float32 = 300
	var height float32 = 100
	var tile int = 0

	if n.Input().GetButton(2) {
		boundingBox := n.GL.BoundingBox()
		mousePosition := n.Input().GetMousePosition()
		mp4 := Vector4{
			X: (mousePosition.X / (boundingBox.Width / 2)) - 1,
			Y: ((boundingBox.Height - mousePosition.Y) / (boundingBox.Height / 2)) - 1,
			Z: 0,
			W: 1,
		}
		mp4 = getProjection().Inverse().MultiplyVector4(mp4)
		width = mp4.X
		height = mp4.Y
	}

	if n.Input().GetKey(n.KeyOne) || n.Input().GetButton(0) {
		tile = 1
	}
	if n.Input().GetKey(n.KeyTwo) {
		tile = 2
	}
	if n.Input().GetKey(n.KeyThree) {
		tile = 3
	}
	if n.Input().GetKey(n.KeyFour) {
		tile = 4
	}
	if n.Input().GetKey(n.KeyFive) {
		tile = 5
	}

	app.slice.Draw(Vector2{0, 0}, Vector2{width, height}, tile)

	/*
		mesh := []Vector3{
			n.NewVector3(0, 0, 0),
			n.NewVector3(x, 0, 0),
			n.NewVector3(0, y, 0),
			n.NewVector3(x, y, 0),
		}

		indecies := []uint16{
			0, 1, 2,
			2, 1, 3,
		}

		n.GL.Clear(n.GlDepthBufferBit)
		app.DrawMesh(mesh, indecies)
	*/

}

func (app *Application) CreateRenderer() {
	indexBuffer = n.GL.CreateBuffer()
	vertexBuffer = n.GL.CreateBuffer()

	a_Vertex = app.shader.GetAttribLocation("a_Vertex")
	u_Projection = app.shader.GetUniformLocation("u_Projection")
}

func (app *Application) DrawMesh(vertices []Vector3, indecies []uint16) {
	// Setup the shader
	n.GL.UseProgram(app.shader.GetProgram())
	n.GL.Enable(n.GlDepthTest)

	// Bind the vertex
	n.GL.BindBuffer(n.GlArrayBuffer, vertexBuffer)                // Bind the buffer to the GlArrayBuffer
	n.GL.BufferData(n.GlArrayBuffer, vertices, n.GlStaticDraw)    // Set teh data to the buffer
	n.GL.VertexAttribPointer(a_Vertex, 3, n.GlFloat, false, 0, 0) // Declare the a_Vertex as the vertex data (while bound)
	n.GL.EnableVertexAttribArray(a_Vertex)                        // Enable the vertex data

	// Bind the indicies
	n.GL.BindBuffer(n.GlElementArrayBuffer, indexBuffer)
	n.GL.BufferData(n.GlElementArrayBuffer, indecies, n.GlStaticDraw)

	// Bind the uniformed
	projection := getProjection()
	n.GL.UniformMatrix4fv(u_Projection, projection)

	// Draw the elements
	n.GL.DrawElements(n.GlTriangles, 6, n.GlUnsignedShort, 0)
}

func getProjection() Matrix {
	//m := n.NewMatrixOrtho(0, 0, , , 400, -400)
	//m = m.Translate(cameraPosition.Negate())

	width := float32(n.GL.Width()) / pixelScale
	height := float32(n.GL.Height()) / pixelScale
	depth := float32(1000)

	m := n.NewMatrixOrtho(0, width, height, 0, depth, -depth)
	m = m.Translate(cameraPosition.Negate())
	m = m.RotateX(0)
	m = m.RotateY(0)
	m = m.RotateZ(0)
	m = m.Scale(n.NewVector3(1, 1, 1))
	return m
}

func transform(vector Vector3, matrix Matrix) Vector3 {
	a := vector.DecomposePointer()
	m := matrix.DecomposePointer()
	var x = a[0]
	var y = a[1]
	var z = a[2]
	var w = m[3]*x + m[7]*y + m[11]*z + m[15]
	if w == 0 {
		w = 1
	}
	return Vector3{
		(m[0]*x + m[4]*y + m[8]*z + m[12]) / w,
		(m[1]*x + m[5]*y + m[9]*z + m[13]) / w,
		(m[2]*x + m[6]*y + m[10]*z + m[14]) / w,
	}
}

func invert(m Matrix) Matrix {
	log.Println("===")
	a := m.Decompose()
	var a00 = a[0]
	var a01 = a[1]
	var a02 = a[2]
	var a03 = a[3]
	var a10 = a[4]
	var a11 = a[5]
	var a12 = a[6]
	var a13 = a[7]
	var a20 = a[8]
	var a21 = a[9]
	var a22 = a[10]
	var a23 = a[11]
	var a30 = a[12]
	var a31 = a[13]
	var a32 = a[14]
	var a33 = a[15]

	var b00 = a00*a11 - a01*a10
	var b01 = a00*a12 - a02*a10
	var b02 = a00*a13 - a03*a10
	var b03 = a01*a12 - a02*a11
	var b04 = a01*a13 - a03*a11
	var b05 = a02*a13 - a03*a12
	var b06 = a20*a31 - a21*a30
	var b07 = a20*a32 - a22*a30
	var b08 = a20*a33 - a23*a30
	var b09 = a21*a32 - a22*a31
	var b10 = a21*a33 - a23*a31
	var b11 = a22*a33 - a23*a32

	log.Print("b00:: ", a00*a11-a01*a10, a00*a11, a01*a10)

	/*    |	0				1			2			3
	------+------------------------------------------------
		0 |	0.9641434 		0 			0 			0
		1 |	0 				0.99999994 	0 			0
		2 |	0 				0 			-1 			-1
		3 |	0 				0 			1.370196 	1.370196
	*/
	var det = (b00 * b11) - (b01 * b10) + (b02 * b09) + (b03 * b08) - (b04 * b07) + (b05 * b06)
	log.Println(a)
	log.Println("b00 * b11", b00*b11, b00, b11)
	log.Println("b01 * b10", b01*b10, b01, b10)
	log.Println("b02 * b09", b02*b09, b02, b09)
	log.Println("b03 * b08", b03*b08, b03, b08)
	log.Println("b04 * b07", b04*b07, b04, b07)
	log.Println("b05 * b06", b05*b06, b05, b06)
	if det == 0 {
		return n.NewMatrix()
	}

	det = 1.0 / det
	return n.Matrix{
		(a11*b11 - a12*b10 + a13*b09) * det,
		(a02*b10 - a01*b11 - a03*b09) * det,
		(a31*b05 - a32*b04 + a33*b03) * det,
		(a22*b04 - a21*b05 - a23*b03) * det,
		(a12*b08 - a10*b11 - a13*b07) * det,
		(a00*b11 - a02*b08 + a03*b07) * det,
		(a32*b02 - a30*b05 - a33*b01) * det,
		(a20*b05 - a22*b02 + a23*b01) * det,
		(a10*b10 - a11*b08 + a13*b06) * det,
		(a01*b08 - a00*b10 - a03*b06) * det,
		(a30*b04 - a31*b02 + a33*b00) * det,
		(a21*b02 - a20*b04 - a23*b00) * det,
		(a11*b07 - a10*b09 - a12*b06) * det,
		(a00*b09 - a01*b07 + a02*b06) * det,
		(a31*b01 - a30*b03 - a32*b00) * det,
		(a20*b03 - a21*b01 + a22*b00) * det,
	}
}

var _debugDiv js.Value
var _hasDebugDiv = false

func drawDebugDiv(position Vector3) {
	if !_hasDebugDiv {
		_hasDebugDiv = true
		document := js.Global().Get("document")
		_debugDiv = document.Call("createElement", "div")
		_debugDiv.Get("classList").Call("add", "spaget-indicator")
		document.Get("body").Call("appendChild", _debugDiv)
		log.Println("shit")
	}

	bounding := n.Canvas().Call("getBoundingClientRect")
	top := position.Y + float32(bounding.Get("top").Float())
	left := position.X + float32(bounding.Get("left").Float())
	_debugDiv.Get("style").Set("top", top)
	_debugDiv.Get("style").Set("left", left)
}

var quadVerts = []Vector3{
	// left column
	Vector3{0, 0, 0},
	Vector3{1, 0, 0},
	Vector3{1, 1, 0},
	Vector3{0, 1, 0},
}

var quadIndicies = []uint16{
	2, 3, 0,
	0, 1, 2,
}

/*

	mousePosition := n.Input().GetMousePosition()
	// Draw nodes
	app.nodeRenderer.Begin()
	app.nodeRenderer.SetSprite(app.tile)
	app.nodeRenderer.Zoom = 0.20
	app.nodeRenderer.SetScale(50)
	//uiMousePosition := app.nodeRenderer.Screen2UISpace(mousePosition)
	app.nodeRenderer.Draw(n.NewRectangle(100, 100, 50, 50), n.Red)
	app.nodeRenderer.Draw(n.NewRectangle(10, 10, 50, 50), n.White)
	app.nodeRenderer.End()

	// Finally render the mouse
	app.spriteRenderer.Begin()
	n.GL.Enable(n.GlDepthTest)
	n.GL.Enable(n.GlBlend)
	n.GL.BlendFunc(n.GlSrcAlpha, n.GlOneMinusSrcAlpha)
	n.GL.DepthFunc(n.GlLess)
	mouseTransform := n.NewTransform2D(mousePosition.Add(cursorHotspot), 0, cursorSize)
	appleTransform := n.NewTransform2D(Vector2{100, 100}, 0, Vector2{10, 10})
	app.spriteRenderer.Draw(app.cursor, Vector2{0, 0}, mouseTransform, n.White)
	app.spriteRenderer.Draw(app.cursor, Vector2{0, 0}, appleTransform, n.Red)
	app.spriteRenderer.End()
*/
