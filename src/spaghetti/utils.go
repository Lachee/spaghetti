package spaghetti

import (
	"syscall/js"
	"unsafe"

	n "github.com/lachee/noodle"
)

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

//TypedArrayToFloat32Slice Converts the JS Float32Array to a slice of float32
func TypedArrayToFloat32Slice(jsValue js.Value) []float32 {
	buff := js.Global().Get("Uint8Array").New(jsValue.Get("buffer"))
	bytes := make([]byte, jsValue.Get("length").Int()*4)
	js.CopyBytesToGo(bytes, buff)
	return *(*[]float32)(unsafe.Pointer(&bytes))
}

//TypedArrayToFloat32Slice Converts the JS Float32Array to a slice of float32
func TypedArrayToUint16Slice(jsValue js.Value) []uint16 {
	buff := js.Global().Get("Uint8Array").New(jsValue.Get("buffer"))
	bytes := make([]byte, jsValue.Get("length").Int()*2)
	js.CopyBytesToGo(bytes, buff)
	return *(*[]uint16)(unsafe.Pointer(&bytes))
}

func checkGLPanic() {
	err := n.GL.Error()
	if err != nil {
		panic(err)
	}
}
