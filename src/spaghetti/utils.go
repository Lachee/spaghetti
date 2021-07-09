package spaghetti

import n "github.com/lachee/noodle"

//GetUIMousePosition gets the position of the mouse on the UI
func GetUIMousePosition() Vector2 {
	return CanvasToUI(n.Input().GetMousePosition())
}

//CanvasToUI converts the canvas pixel values to screen space
func CanvasToUI(point Vector2) Vector2 {
	boundingBox := n.GL.BoundingBox()
	projected := getProjection().Inverse().MultiplyVector4(Vector4{
		X: (point.X / (boundingBox.Width / 2)) - 1,
		Y: ((boundingBox.Height - point.Y) / (boundingBox.Height / 2)) - 1,
		Z: 0,
		W: 1,
	})
	return n.NewVector2(projected.X, projected.Y)
}

func checkGLPanic() {
	err := n.GL.Error()
	if err != nil {
		panic(err)
	}
}


