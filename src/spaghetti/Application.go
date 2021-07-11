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

var (
	cameraPosition Vector3
	pixelScale     float32

	hue float32
)

//Application handles the game. Put your variables in here
type Application struct {
	shader *n.Shader
	bg     *Background

	Renderer *Renderer

	boxWindow SliceWindow
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
	renderer, err := NewRenderer()
	app.Renderer = renderer
	if err != nil {
		return false
	}

	// Create the texture
	boxImage, err := LoadResourceImage("resource://textures/slice_atlas.png")
	if err != nil {
		return false
	}

	// Create the window
	boxTexture := boxImage.CreateTexture()
	boxTexture.SetFilter(n.TextureFilterNearest)
	app.boxWindow = SliceWindow{
		texture: boxTexture,
		top:     5,
		bottom:  5,
		left:    5,
		right:   5,
		cols:    6,
		rows:    1,
	}

	// Prepare font resource
	fontResourceResult := <-FetchResource("resource://font/LobsterTwo-Regular.ttf")
	if fontResourceResult.Error != nil {
		n.Error("Failed to load the font", fontResourceResult.Error)
		return false
	}

	// Get the polygons
	fontPaths := fontResourceResult.Data.Invoke("Hello World", 12)
	log.Println(fontPaths)

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

	// Colour shifting
	hue += 90 * dt

	//app.Renderer.Update()
	//log.Println("Camera Position", cameraPosition, axis)
}

//Render draws the frame
func (app *Application) Render() {

	n.DebugDraw = true

	// Clear the canvas
	n.GL.ClearColor(n.White)
	n.GL.Clear(n.GlColorBufferBit | n.GlDepthBufferBit)
	//app.bg.Draw()
	//n.GL.Clear(n.GlDepthBufferBit)

	app.Renderer.DrawBox(n.NewRectangle(10, 10, 150, 150), Point{0, 0}, app.boxWindow)
	app.Renderer.DrawBox(n.NewRectangle(30, 30, 150, 150), Point{1, 0}, app.boxWindow)
	app.Renderer.DrawBox(n.NewRectangle(50, 50, 150, 150), Point{2, 0}, app.boxWindow)

	mouse := GetUIMousePosition()
	//app.Renderer.DrawBox(n.NewRectangle(70, 70, mouse.X-70, mouse.Y-70), Point{3, 0}, app.boxWindow)
	app.Renderer.DrawRectangle(n.NewRectangle(70, 70, mouse.X-70, mouse.Y-70), n.NewColorFromHSV(n.NewVector3(hue, 1, 1)))
	// Finally render
	app.Renderer.Render()

	/*
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
	*/

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
