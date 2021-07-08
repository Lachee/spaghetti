package spaghetti

import (
	n "github.com/lachee/noodle"
)

type Button struct {
	rectangle n.Rectangle
	text      string

	normalStyle, hoverStyle, pressedStyle int
}

//getStyle gets the current style state
func (btn Button) getStyle() int {
	if btn.rectangle.Contains(GetMousePosition()) {
		if n.Input().GetButton(0) {
			return btn.pressedStyle
		}
		return btn.hoverStyle
	}
	return btn.normalStyle
}

type UI struct {
	matSlice *NineSliceMaterial

	buttons []Button

	backgroundAtlas *n.Texture
	style           NineSliceCut
}

//NewUI creates a new UI instance
func NewUI() (*UI, error) {
	ui := &UI{}

	sliceMaterial, sliceError := NewNineSliceMaterial()
	ui.matSlice = sliceMaterial
	if sliceError != nil {
		return nil, sliceError
	}

	backgroundImage, backgroundError := LoadResourceImage("resource://textures/slice_atlas.png")
	if backgroundError != nil {
		return nil, backgroundError
	}

	ui.backgroundAtlas = backgroundImage.CreateTexture()
	ui.backgroundAtlas.SetFilter(n.TextureFilterNearest)
	ui.style = NineSliceCut{
		texture: ui.backgroundAtlas,
		cols:    6, rows: 1,
		top: 5, bottom: 5,
		left: 5, right: 5,
	}

	return ui, nil
}

func (ui *UI) AddButton(button Button) {
	ui.buttons = append(ui.buttons, button)
}

func (ui *UI) Draw() {
	ui.drawBackgrounds()
}

func (ui *UI) drawBackgrounds() {

	o := float32(0)
	i := float32(1)

	verticies := make([]float32, 0)
	indicies := make([]uint16, 0)

	for k, btn := range ui.buttons {

		size := btn.rectangle.Size()
		style := btn.getStyle()
		tileX := float32(style % ui.style.cols)
		tileY := float32(style / ui.style.cols)

		verticies = append(verticies,
			// Vertex, UV, Size, Offset
			btn.rectangle.X, btn.rectangle.Y, o, // 0 0 0
			o, o,
			size.X, size.Y,
			tileX, tileY,

			btn.rectangle.X+btn.rectangle.Width, btn.rectangle.Y, o, // 1 0 0
			i, o,
			size.X, size.Y,
			tileX, tileY,

			btn.rectangle.X, btn.rectangle.Y+btn.rectangle.Height, o, // 0 1 0
			o, i,
			size.X, size.Y,
			tileX, tileY,

			btn.rectangle.X+btn.rectangle.Width, btn.rectangle.Y+btn.rectangle.Height, o, // 1 1 0
			i, i,
			size.X, size.Y,
			tileX, tileY,
		)

		index := uint16(k * 4)
		indicies = append(indicies,
			index+0, index+1, index+2,
			index+2, index+1, index+3,
		)
	}

	ui.matSlice.SetCut(ui.style)
	ui.matSlice.Draw(verticies, indicies)
}
