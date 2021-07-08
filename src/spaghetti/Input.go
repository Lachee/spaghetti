package spaghetti

import n "github.com/lachee/noodle"

func GetMousePosition() Vector2 {
	boundingBox := n.GL.BoundingBox()
	mousePosition := n.Input().GetMousePosition()
	mp4 := Vector4{
		X: (mousePosition.X / (boundingBox.Width / 2)) - 1,
		Y: ((boundingBox.Height - mousePosition.Y) / (boundingBox.Height / 2)) - 1,
		Z: 0,
		W: 1,
	}
	mp4 = getProjection().Inverse().MultiplyVector4(mp4)
	return n.NewVector2(mp4.X, mp4.Y)
}
