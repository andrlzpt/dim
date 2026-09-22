package renderer

import (
	"github.com/andrlzpt/dim/internal/ascii"
	"github.com/go-gl/gl/v3.2-compatibility/gl"
	"golang.org/x/image/font/basicfont"
)

// RenderText draws shaded glyphs over existing layers without clearing.
// Call BeginFrame first. Use 7x13 cells or integer multiples such as 14x26.
func RenderText(buffer *ascii.Buffer, config Config) {
	gl.Begin(gl.QUADS)
	for y := 0; y < buffer.Height; y++ {
		for x := 0; x < buffer.Width; x++ {
			cell := buffer.Cells[y*buffer.Width+x]
			if cell.Rune == 0 || cell.Rune == ' ' || cell.Brightness <= 0 {
				continue
			}

			gl.Color3f(
				cell.Brightness,
				cell.Brightness,
				cell.Brightness,
			)
			drawTextGlyph(x, y, cell.Rune, config, nil)
		}
	}
	gl.End()
}

// RenderTextLine draws a centered label at pixel coordinates.
// quarterTurns: 0 horizontal, 1 clockwise, -1 counterclockwise.
func RenderTextLine(
	label string,
	centerX, centerY float32,
	quarterTurns int,
	config Config,
) {
	runes := []rune(label)
	width := float32(len(runes) * config.CellWidth)
	height := float32(config.CellHeight)

	transform := func(x, y float32) (float32, float32) {
		x -= width / 2
		y -= height / 2

		switch quarterTurns {
		case 1:
			x, y = -y, x
		case -1:
			x, y = y, -x
		}

		return centerX + x, centerY + y
	}

	gl.Color3f(1, 1, 1)
	gl.Begin(gl.QUADS)
	for x, r := range runes {
		if r != ' ' {
			drawTextGlyph(x, 0, r, config, transform)
		}
	}
	gl.End()
}

func drawTextGlyph(cellX, cellY int, r rune, config Config, transform func(float32, float32) (float32, float32)) {
	face := basicfont.Face7x13

	glyphWidth := config.GlyphWidth
	if glyphWidth <= 0 {
		glyphWidth = config.CellWidth
	}

	glyphHeight := config.GlyphHeight
	if glyphHeight <= 0 {
		glyphHeight = config.CellHeight
	}

	scaleX := float32(glyphWidth) / float32(face.Advance)
	scaleY := float32(glyphHeight) / float32(face.Height)
	originX := float32(cellX * config.CellWidth)
	originY := float32(cellY * config.CellHeight)

	vertex := func(x, y float32) {
		if transform != nil {
			x, y = transform(x, y)
		}
		gl.Vertex2f(
			pixelToClipX(x, config),
			pixelToClipY(y, config),
		)
	}

	for _, pixel := range cachedGlyphPixels(r) {
		x := originX + pixel.X*scaleX
		y := originY + pixel.Y*scaleY

		vertex(x, y)
		vertex(x+scaleX, y)
		vertex(x+scaleX, y+scaleY)
		vertex(x, y+scaleY)
	}
}
