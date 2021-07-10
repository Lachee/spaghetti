package spaghetti

import n "github.com/lachee/noodle"

type Renderer struct {
	sliceRender *SliceRender
	commands    []RenderCommand
}

//NewRenderer creates a new UI instance
func NewRenderer() (*Renderer, error) {
	renderer := &Renderer{}

	sliceRender, err := NewSliceMaterial()
	renderer.sliceRender = sliceRender
	if err != nil {
		n.Error("Failed to load the renderer", err)
		return nil, err
	}

	return renderer, nil
}

func (renderer *Renderer) DrawRectangle(rectangle Rectangle, color Color) {
	// https://www.shadertoy.com/view/3tj3Dm
}

func (renderer *Renderer) DrawBox(rectangle Rectangle, tile Point, window SliceWindow) {
	cmd := &boxCommand{
		rectangle: rectangle,
		tile:      tile,
		window:    window,
	}

	renderer.push(cmd)
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

type rectangleCommand struct {
}

func (cmd *rectangleCommand) Render(ctx RenderContext) {
}

//boxCommand
type boxCommand struct {
	rectangle Rectangle
	tile      Point
	window    SliceWindow
}

func (cmd *boxCommand) Render(ctx RenderContext) {
	ctx.renderer.sliceRender.Draw(cmd.rectangle, cmd.tile, cmd.window)
}
