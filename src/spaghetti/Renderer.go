package spaghetti

import n "github.com/lachee/noodle"

type Render interface {
	//Render the buffered elements
	Render()
}

type Renderer struct {
	sliceRender *SliceRender
	flatRender  *FlatRender
	commands    []RenderCommand
}

//NewRenderer creates a new UI instance
func NewRenderer() (*Renderer, error) {
	renderer := &Renderer{}

	sliceRender, err := NewSliceRender()
	renderer.sliceRender = sliceRender
	if err != nil {
		n.Error("Failed to load the slice render", err)
		return nil, err
	}

	flatRender, err := NewFlatRender()
	renderer.flatRender = flatRender
	if err != nil {
		n.Error("Failed to load the flat render", err)
		return nil, err
	}

	return renderer, nil
}

func (renderer *Renderer) DrawRectangle(rectangle Rectangle, color Color) {
	// TODO: Implement rounded corners like so - https://www.shadertoy.com/view/3tj3Dm
	renderer.push(&rectangleCommand{
		rectangle: rectangle,
		color:     color,
	})
}

func (renderer *Renderer) DrawBox(rectangle Rectangle, tile Point, window SliceWindow) {
	renderer.push(&boxCommand{
		rectangle: rectangle,
		tile:      tile,
		window:    window,
	})
}

//push a command to the render queue
func (renderer *Renderer) push(command RenderCommand) {
	renderer.commands = append(renderer.commands, command)
}

// Render the current queue
func (renderer *Renderer) Render() {
	// process the queue
	for i, cmd := range renderer.commands {
		cmd.Render(RenderContext{
			depth:    float32(i),
			renderer: renderer,
		})
	}

	// render any slices that were still being buffered
	renderer.sliceRender.Render()
	renderer.flatRender.Render()

	// clear the queue
	renderer.commands = renderer.commands[:0]
}

type RenderContext struct {
	depth    float32
	renderer *Renderer
}

type RenderCommand interface {
	Render(ctx RenderContext)
}

//rectangleCommand draws a flat rectangle on the screen
type rectangleCommand struct {
	rectangle Rectangle
	color     Color
}

func (cmd *rectangleCommand) Render(ctx RenderContext) {
	ctx.renderer.flatRender.DrawRectangle(cmd.rectangle, cmd.color)
}

//boxCommand drwas a box on the screen
type boxCommand struct {
	rectangle Rectangle
	tile      Point
	window    SliceWindow
}

func (cmd *boxCommand) Render(ctx RenderContext) {
	ctx.renderer.sliceRender.Draw(cmd.rectangle, cmd.tile, cmd.window)
}
